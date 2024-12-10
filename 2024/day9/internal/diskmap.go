package internal

import (
	"fmt"
	"math/big"
	"sort"
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

func (dm *DiskMap) String() string {
	// Find the maximum block index to determine the range of blocks.
	maxBlockIndex := -1
	for blockIndex := range dm.blockToFile {
		if blockIndex > maxBlockIndex {
			maxBlockIndex = blockIndex
		}
	}

	// Build the string representation.
	var result string
	result += "["
	for i := 0; i <= maxBlockIndex; i++ {
		if fileID, exists := dm.blockToFile[i]; exists {
			result += fmt.Sprintf("%d", fileID)
		} else {
			result += "."
		}
	}
	result += "]"

	return result
}


// Compact reorganizes the DiskMap by moving the highest block ID to the smallest available empty block ID.
// It does not attempt to keep file blocks together.
func (dm *DiskMap) Compact() error {
	if len(dm.blockToFile) == 0 {
		// No blocks to compact.
		return nil
	}

	// Step 1: Identify all occupied blocks and determine the maximum block ID.
	maxBlockID := -1
	for block := range dm.blockToFile {
		if block > maxBlockID {
			maxBlockID = block
		}
	}

	// Step 2: Identify all empty blocks (assuming blocks start at 0 and are contiguous up to maxBlockID).
	emptyBlocks := make([]int, 0)
	occupied := make(map[int]bool)
	for block := range dm.blockToFile {
		occupied[block] = true
	}

	for block := 0; block < maxBlockID; block++ {
		if !occupied[block] {
			emptyBlocks = append(emptyBlocks, block)
		}
	}

	if len(emptyBlocks) == 0 {
		// No empty blocks to compact.
		return nil
	}

	// Step 3: Sort empty blocks in ascending order and occupied blocks in descending order.
	sort.Ints(emptyBlocks)
	occupiedBlocks := make([]int, 0)
	for block := range dm.blockToFile {
		occupiedBlocks = append(occupiedBlocks, block)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(occupiedBlocks)))

	// Step 4: Iterate and move blocks.
	for _, sourceBlock := range occupiedBlocks {
		if len(emptyBlocks) == 0 {
			break // No more empty blocks to fill.
		}

		// If the source block is already in a lower position than the smallest empty block, skip.
		if sourceBlock <= emptyBlocks[0] {
			continue
		}

		// Move the source block to the smallest empty block.
		targetBlock := emptyBlocks[0]

		fileID := dm.blockToFile[sourceBlock]

		// Update fileBlockIndex for the file.
		blocks := dm.fileBlockIndex[fileID]
		for i, b := range blocks {
			if b == sourceBlock {
				dm.fileBlockIndex[fileID][i] = targetBlock
				break
			}
		}

		// Update blockToFile mapping.
		delete(dm.blockToFile, sourceBlock)
		dm.blockToFile[targetBlock] = fileID

		// Update occupied and empty blocks lists.
		emptyBlocks = emptyBlocks[1:]
		occupied[sourceBlock] = false
		occupied[targetBlock] = true
	}

	// Step 5: Update the checksum after compaction.
	dm.UpdateChecksum()

	return nil
}
