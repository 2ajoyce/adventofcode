from pathlib import Path
import sys
import shutil


def eprint(*args, **kwargs):
    print(*args, **kwargs, file=sys.stderr)


def main():
    if len(sys.argv) != 2:
        eprint("Usage: python scripts/start-day.py <day>")
        sys.exit(2)

    day = sys.argv[1]
    if not day.isdigit() or not (1 <= int(day) <= 31):
        eprint("Error: <day> must be an integer between 1 and 31.")
        sys.exit(2)

    # locate project root relative to this script
    script_path = Path(__file__).resolve()
    project_root = script_path.parent.parent

    template_dir = project_root / "template"
    if not template_dir.exists():
        eprint(f"Error: template directory not found at {template_dir}")
        sys.exit(1)

    day_dir = project_root / str(day)

    if day_dir.exists():
        print(
            f"Directory '{day}' already exists. Will not overwrite existing files; copying missing files only."
        )
    else:
        day_dir.mkdir(parents=True, exist_ok=True)
        print(f"Created directory '{day}'")

    copied = 0
    skipped = 0

    for src in template_dir.rglob("*"):
        rel = src.relative_to(template_dir)
        dest = day_dir / rel

        if src.is_dir():
            dest.mkdir(parents=True, exist_ok=True)
            continue

        if dest.exists():
            skipped += 1
            print(f"Skipping existing file: {dest}")
            continue

        # ensure parent directory exists
        dest.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(src, dest)
        copied += 1
        print(f"Copied: {dest}")

        # If we copied a go.mod, update its module path so the last
        # path segment matches the target day directory (e.g. replace
        # trailing 'template' with '1'). This keeps module names
        # consistent per-day.
        if dest.name == "go.mod":
            try:
                text = dest.read_text(encoding="utf-8")
            except Exception:
                # If we can't read the file for some reason, continue.
                print(f"Warning: couldn't read {dest} to patch module line.")
            else:
                lines = text.splitlines()
                for i, line in enumerate(lines):
                    if line.strip().startswith("module "):
                        mod = line.strip()[len("module ") :].strip()
                        parts = mod.split("/")
                        if parts:
                            parts[-1] = day
                            newmod = "/".join(parts)
                            lines[i] = "module " + newmod
                        break
                # write back; ensure newline at EOF
                dest.write_text("\n".join(lines) + "\n", encoding="utf-8")
                print(f"Patched module path in: {dest}")

    print(f"Done. Copied {copied} file(s), skipped {skipped} existing file(s).")


if __name__ == "__main__":
    main()
