# engineSync

Application to sync EngineDJ (Desktop) Library Backups to drives.

Works (developed and tested) on macOS Ventura 🎉

This tool was created because macOS Ventura changed something with EXFAT/FAT handling which in turn broke EngineDJ. 
This tool simply copies the necessary files and makes adjustments to the file paths in the engine library, it doesn't check for a specific partition format and so is not affected by that bug.

This tool is not supported by Denon/InMusic; Please read the disclaimer at the bottom of the page.

## Features

- Syncs entire library to a given drive/directory
- Syncs a SoundSwitch project to the target drive
- Fixes DB entries for track paths according to what's been synced
- Works on macOS Ventura (and probably every OS golang runs on)

### What it can't do

- Sync stuff back into your Engine DJ Library
- Recognize existing libraries 
  - it overwrites everything without checking
- Work with non-Backup folders (ie. a regular library)
  - they contain "invalid" paths to the tracks, the backup contains the proper working paths

## Usage

1. Download the tool (https://github.com/lfuelling/engineSync/releases/latest)
2. Organize your Music using engineDJ (which runs good enough for that purpose on macOS Ventura)
3. (OPTIONAL) Organize your light show using SoundSwitch (also still running on Ventura)
4. Close engineDJ, let it create a **Library Backup** (this is important)
5. Attach your drive, make sure it's formatted to EXFAT or FAT!
   - **I have bricked my drive previously with the bugged engineDJ, so my drive was freshly formatted**
   - **I have not tested this tool with an existing library on the drive**
   - **it might not work and will definitely leave duplicate files on the drive if it already contains a library**
6. Follow the instructions in the tool window and wait for the sync to finish
   1. Select the "Engine Library Backup" folder
   2. (Optional) select the SoundSwitch project
   3. Select the target drive
   4. Click "Start Sync" and wait
7. **Properly Eject your drive!**
8. Insert your drive into any engineDJ player and enjoy.

## Development

(macOS only so far, should run anywhere golang does though)

1. Clone the repo
2. Install dependencies:
   - Download or install golang
3. Build:
   - `go build engineSync`
4. Run:
   - `./engineSync`

## Disclaimer

This software is not endorsed, supported by or in any way affiliated with DenonDJ and/or InMusic Brands. All the respective Copyrights belong to them. None of their software and/or hardware products was disassembled or in any way used against the license requirements to develop this tool.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. 
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
