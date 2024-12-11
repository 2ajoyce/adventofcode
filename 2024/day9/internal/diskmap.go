package internal

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
)

// DiskMap represents the mapping of files to disk blocks.
type DiskMap struct {
	// Maps a file ID to its size in blocks.
	fileSize map[int]int

	// Maps a file ID to a slice of block indices that make up the file.
	fileBlockIndex map[int][]int

	// Maps a block index to its corresponding file ID for efficient lookup.
	blockToFile map[int]int

	// Checksum representing the current state.
	checksum big.Int
}

// NewDiskMap initializes and returns a new DiskMap.
func NewDiskMap() *DiskMap {
	return &DiskMap{
		fileSize:       make(map[int]int),
		fileBlockIndex: make(map[int][]int),
		blockToFile:    make(map[int]int),
	}
}

// AddFile adds a file with the given fileID and associated blocks to the DiskMap.
func (dm *DiskMap) AddFile(fileID int, blocks []int) error {
	if _, exists := dm.fileSize[fileID]; exists {
		return fmt.Errorf("file ID %d already exists", fileID)
	}

	dm.fileSize[fileID] = len(blocks)
	dm.fileBlockIndex[fileID] = make([]int, len(blocks))
	copy(dm.fileBlockIndex[fileID], blocks)

	// Update the reverse mapping.
	for _, block := range blocks {
		if _, occupied := dm.blockToFile[block]; occupied {
			return fmt.Errorf("block %d is already occupied by file ID %d", block, dm.blockToFile[block])
		}
		dm.blockToFile[block] = fileID
	}

	return nil
}

// GetFileId returns the file ID associated with the given block.
// Returns an error if the block is not found.
func (dm *DiskMap) GetFileId(block int) (int, error) {
	if fileID, exists := dm.blockToFile[block]; exists {
		return fileID, nil
	}
	return -1, fmt.Errorf("block %d not found in any file", block)
}

// GetFileSize returns the size of the file in blocks.
func (dm *DiskMap) GetFileSize(fileID int) (int, error) {
	size, exists := dm.fileSize[fileID]
	if !exists {
		return -1, fmt.Errorf("file ID %d not found", fileID)
	}
	return size, nil
}

// GetFileBlocks returns the list of block indices that make up the file.
func (dm *DiskMap) GetFileBlocks(fileID int) ([]int, error) {
	blocks, exists := dm.fileBlockIndex[fileID]
	if !exists {
		return nil, fmt.Errorf("file ID %d not found", fileID)
	}
	return blocks, nil
}

// GetChecksum returns the current checksum.
func (dm *DiskMap) GetChecksum() *big.Int {
	return new(big.Int).Set(&dm.checksum)
}

// UpdateChecksum recalculates and updates the checksum based on current block assignments.
func (dm *DiskMap) UpdateChecksum() *big.Int {
	dm.checksum.SetInt64(0)

	for fileID, blocks := range dm.fileBlockIndex {
		for _, block := range blocks {
			fileIDBig := big.NewInt(int64(fileID))
			blockBig := big.NewInt(int64(block))
			product := new(big.Int).Mul(fileIDBig, blockBig)
			dm.checksum.Add(&dm.checksum, product)
		}
	}

	return dm.GetChecksum()
}

// String returns a string representation of the DiskMap.
// Occupied blocks are represented by their file IDs, and empty blocks by '.'.
// Example: [1..2.1.2.1.3333]
func (dm *DiskMap) String() string {
	if len(dm.blockToFile) == 0 {
		return "[]"
	}

	// Determine the range of blocks.
	var maxBlock int
	for block := range dm.blockToFile {
		if block > maxBlock {
			maxBlock = block
		}
	}

	// Use a strings.Builder for efficient string concatenation.
	var sb strings.Builder
	sb.WriteString("[")

	for block := 0; block <= maxBlock; block++ {
		if fileID, exists := dm.blockToFile[block]; exists {
			sb.WriteString(fmt.Sprintf("%d", fileID))
		} else {
			sb.WriteString(".")
		}
	}

	sb.WriteString("]")

	return sb.String()
}

// Compact reorganizes the DiskMap by moving whole files to the leftmost span of free blocks that can fit the file.
// It processes files in order of decreasing file ID and attempts to move each file exactly once.
func (dm *DiskMap) Compact() error {
	if len(dm.blockToFile) == 0 {
		// No blocks to compact.
		return nil
	}

	// Step 1: Sort files in order of decreasing file ID.
	fileIDs := make([]int, 0, len(dm.fileSize))
	for fileID := range dm.fileSize {
		fileIDs = append(fileIDs, fileID)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(fileIDs)))

	// Step 2: Attempt to move each file once.
	for _, fileID := range fileIDs {
		// Moving Step3 and Step4 inside the loop made this much, much less performant
		// To increase performance, the bulk search should be done once outside the loop
		// with individual updates inside the loop to the newly empty and full spans.

		// Step 3: Identify all free blocks and determine free spans.
		freeSpans := dm.findFreeSpans()

		// Step 4: Sort free spans by starting block index (leftmost first).
		sort.Slice(freeSpans, func(i, j int) bool {
			return freeSpans[i].start < freeSpans[j].start
		})

		fileSize := dm.fileSize[fileID]
		if fileSize <= 0 {
			continue // Skip empty files
		}

		// Find the first free span that can fit the entire file.
		var suitableSpan *Span
		for _, span := range freeSpans {
			if span.length >= fileSize {
				suitableSpan = &span
				break
			}
		}

		// Find the smallest block ID in the current file's blocks.
		currentBlocks := dm.fileBlockIndex[fileID]
		smallestBlockId := currentBlocks[0]
		for _, blockId := range currentBlocks {
			if blockId < smallestBlockId {
				smallestBlockId = blockId
			}
		}
		// IF no suitable span is found
		// OR it starts at a block ID larger than the smallest current block ID
		if suitableSpan == nil || suitableSpan.start >= smallestBlockId {
			// No suitable span found; do not move the file.
			continue
		}

		// Determine the target blocks for the file.
		targetStart := suitableSpan.start
		targetBlocks := make([]int, fileSize)
		for i := 0; i < fileSize; i++ {
			targetBlocks[i] = targetStart + i
		}

		// Move the file to the target blocks.
		err := dm.moveFile(fileID, targetBlocks)
		if err != nil {
			return fmt.Errorf("failed to move file %d: %v", fileID, err)
		}
	}

	// Step 5: Update the checksum after compaction.
	dm.UpdateChecksum()

	return nil
}

// Span represents a contiguous range of free blocks.
type Span struct {
	start  int
	length int
}

// findFreeSpans identifies all contiguous free block spans in the DiskMap.
func (dm *DiskMap) findFreeSpans() []Span {
	// Determine the range of blocks.
	var maxBlock int
	for block := range dm.blockToFile {
		if block > maxBlock {
			maxBlock = block
		}
	}

	// Iterate through blocks to find free spans.
	var spans []Span
	currentStart := -1
	currentLength := 0

	for block := 0; block <= maxBlock; block++ {
		if _, occupied := dm.blockToFile[block]; !occupied {
			if currentStart == -1 {
				currentStart = block
				currentLength = 1
			} else {
				currentLength++
			}
		} else {
			if currentStart != -1 {
				spans = append(spans, Span{start: currentStart, length: currentLength})
				currentStart = -1
				currentLength = 0
			}
		}
	}

	// Append the last span if it ends at maxBlock.
	if currentStart != -1 {
		spans = append(spans, Span{start: currentStart, length: currentLength})
	}

	return spans
}

// moveFile moves a file to the specified target blocks.
// It updates fileBlockIndex and blockToFile accordingly.
func (dm *DiskMap) moveFile(fileID int, targetBlocks []int) error {
	// Get current blocks of the file.
	currentBlocks, exists := dm.fileBlockIndex[fileID]
	if !exists {
		return fmt.Errorf("file ID %d not found", fileID)
	}

	// Check if target blocks are free.
	for _, block := range targetBlocks {
		if _, occupied := dm.blockToFile[block]; occupied {
			return fmt.Errorf("target block %d is already occupied", block)
		}
	}

	// Remove current block assignments.
	for _, block := range currentBlocks {
		delete(dm.blockToFile, block)
	}

	// Assign new blocks.
	for _, block := range targetBlocks {
		dm.blockToFile[block] = fileID
	}

	// Update fileBlockIndex.
	dm.fileBlockIndex[fileID] = targetBlocks

	return nil
}
