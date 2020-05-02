# i3wm Workspace window switcher

Window switcher for i3wm workspaces.
This is my attempt to learn go, and my first go program

### Prerequisites

Require `rofi` or `dmenu` to be installed on the syetem


### Installing

```
git clone patrykwegrzyni3wm-workspace-window-switcher
```


```
cd i3wm-workspace-window-switcher
```

if using go

```
go install
```
or
```
cp ./bin/windowswitcher /usr/local/bin
```
## Examples

Dmenu
```
bindsym $mod+Shift+w exec /windowswitcher -nb '#282C34' -sf '#282C34' -sb '#61AEEE' -h 25
```

Rofi
```
bindsym $mod+Shift+w exec /windowswitcher -mode "rofi" -dmenu
```


## Authors

* **Patryk Wegrzyn** - *Initial work* - [PurpleBooth](https://github.com/patrykwegrzyn)


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

