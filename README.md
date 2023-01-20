# HorizonXI Linux Installer

Welcome to the Unofficial [HorizonXI](https://horizonxi.com) Linux Installer

This repo contains an installer HorizonXI on Linux.

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

### Step 1. Download HorizonXI data via bittorrent (HorizonXI.zip)

https://github.com/sheik/horizonxi-linux/blob/4c455018030d6bffe6f7cc7db273355616ffa8f7/cmd/horizonxi-installer/main.go#L30
