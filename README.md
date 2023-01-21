# HorizonXI Linux Installer

Welcome to the Unofficial [HorizonXI](https://horizonxi.com) Linux Installer

## How to Install

[Download the latest HorizonXI Linux Installer here](https://github.com/sheik/horizonxi-linux/releases/download/v0.0.4/horizonxi-installer)

Once downloaded, you need to mark it as executable. Open a terminal and run the following command:

```bash
chmod +x Downloads/horizonxi-installer
```

Finally, you can run the installer from the terminal as well:

```bash
./Downloads/horizonxi-installer
```

Once it is done installing, you can run the game from the terminal with the following command:

```bash
$HOME/HorizonXI/horizonxi
```

## Provide your own data file

Optionally, if you already have HorizonXI.zip downloaded via a different method, you can supply it as an argument to the
installer so that you do not need to download it again:

```bash
./horizonxi-installer -d /path/to/HorizonXI.zip
```

## Troubleshooting

If you run into a problem while installing, you can try removing the install directory and running the installer again:

```bash
rm -rf $HOME/HorizonXI
```

Then:

```bash
./Downloads/horizonxi-installer
```

## Installer Details

I thought it might be helpful to explain exactly what the installer is doing so that steps can be replicated by hand if necessary.

#### 1. Download HorizonXI

The HorizonXI data files (HorizonXI.zip) are downlaoded via bittorrent. This is the current magnet link:

https://github.com/sheik/horizonxi-linux/blob/4c455018030d6bffe6f7cc7db273355616ffa8f7/cmd/horizonxi-installer/main.go#L30-L32

#### 2. Verify Files

After download (or if files are supplied by user), the data files are run through an integrity check using SHA256.

#### Step 3. Unzip HorizonXI.zip into the install directory (default is ~/HorizonXI)

#### Step 4. Install GE-Wine

[https://github.com/GloriousEggroll/wine-ge-custom](wine-ge-custom) is a version of wine that is built from the most
recent bleeding-edge proton experimental wine repo.

#### Step 5. Install DgVoodoo2 (enables graphics card)

[DgVoodoo2](http://dege.freeweb.hu/dgVoodoo2/dgVoodoo2/) is a translation layer that translates older
DirectX/Direct3d/glide APIs into newer versions of those APIs.

This will allow HorizonXI to run on your graphics card instead of the CPU.

#### Step 6. Install DXVK (Direct3D 9-11 -> Vulkan)

[DXVK](https://github.com/doitsujin/dxvk) is another translation layer that translates Direct3D API calls to Vulkan API calls.

This will allow HorizonXI to run better on Linux.

#### Step 7. Patch HorizonXI bootloader to support large address spaces

https://github.com/sheik/horizonxi-linux/blob/4c455018030d6bffe6f7cc7db273355616ffa8f7/cmd/horizonxi-installer/main.go#L264

This will prevent blackscreens when zoning in-game.

#### Step 8. Enable 60 FPS by default

This sets the FPS divisor to 1, which allows a max of 60FPS by default instead of 30FPS.

#### Step 9. Install launcher script



