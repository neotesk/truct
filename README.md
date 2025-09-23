<h1>Truct<img align="left" width="42" height="42" alt="logo" src="https://github.com/user-attachments/assets/1a8052d1-924c-4799-9037-8a7c6ac5fc68" /></h1>

Truct[^1] is a pretty minimal workflow runner, allowing you to store your tasks inside one single
Truct file (generally stored as `truct.yaml`) which is a YAML[^2] file, this way it will be readable
for humans! Compared to other systems like Make[^3], Truct aims to be simple and beginner-friendly
for small projects.

### Installation (Manual)
You can install Truct manually through the [Releases](https://github.com/neotesk/truct/releases)
section. Currently there are builds only for *Nix operating systems (Linux[^4], OpenBSD[^7], macOS[^6] etc)
and Windows[^5].

### Installation (Automatic for Linux/Unix systems)
```
bash <(curl -sSL https://raw.githubusercontent.com/neotesk/truct/main/docs/install.sh)
```

### Installation (Automatic for Termux)
```
bash <(curl -sSL https://raw.githubusercontent.com/neotesk/truct/main/docs/termux-install.sh)
```

### Usage
You can start with the help command like so:
```
truct help
```
After writing your workflow file, you can run workflows with this command:
```
truct do
```

### Why does this exist?
This exists because I like making small projects that will make my job easier and I don't want to
adapt to many many other systems on the current market, so I like to combine my favorite parts of
these systems into one single unit, thus many of my projects have born into existence. Truct is
one of them since I only needed a simple workflow runner that does basic work and nothing else.
For more information, please visit the [Truct wiki](https://github.com/neotesk/truct/wiki)

[^1]: Truct comes from "construct" in English.
[^2]: [YAML](https://en.wikipedia.org/wiki/YAML) is a human-readable markup language, stands for "Yet Another Markup Language"
[^3]: [Make](https://en.wikipedia.org/wiki/Make_(software)) is a command-line interface software tool that performs actions ordered by configured dependencies as defined in a configuration file called a makefile
[^4]: [Linux](https://en.wikipedia.org/wiki/Linux) is an Unix-like Operating system.
[^5]: [Windows](https://en.wikipedia.org/wiki/Microsoft_Windows) is a popular computer operating system used world-wide.
[^6]: [macOS](https://en.wikipedia.org/wiki/MacOS) is a popular computer operating system used in Apple's Mac and Macbook computers.
[^7]: [OpenBSD](https://en.wikipedia.org/wiki/OpenBSD) is a security-focused, free software, Unix-like operating system based on the Berkeley Software Distribution (BSD).