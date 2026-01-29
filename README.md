# motion-morgue

A temporary holding place for recovering the party assembly records before they can be properly archived.

> **Note:** This is a "vibe-coded" tool. The focus is on collecting scattered data, not on building production-quality software. No security considerations, no tests, no error handling beyond the basics. Use at your own risk - and ideally only locally.

## The Problem

Years of assembly documentation - protocols, motions, amendments, election results - exist scattered across various formats, incomplete folders, and personal archives.
Some records are missing entirely.
There is no single source of truth.

## The Purpose

This tool collects and organizes what remains:

- **Assemblies** with their protocols
- **Motions** and their original documents
- **Amendments** to those motions

Each record can be added even when incomplete.
Missing PDFs are marked as placeholders to be filled in later.

## The Goal

Once recovered and validated, this data will be migrated into a proper decision database.

## Quick Start

```bash
# Build
go build -o motion-morgue .

# Set the association (required)
export MOTION_MORGUE_ASSOC="my-party"

# Add an assembly
./motion-morgue add assembly --title "BPT 2024-1" --start 2024-03-15 --end 2024-03-17

# Add a motion
./motion-morgue add motion --assembly 1 --sort "A001" --title "Satzungsänderung"

# Add an amendment
./motion-morgue add amendment --motion 1 --sort "Ä001"

# Import a PDF
./motion-morgue import assembly 1 ./protokoll.pdf

# List everything
./motion-morgue list
```

Data is stored in `~/.motion-morgue/<MOTION_MORGUE_ASSOC>/`.
