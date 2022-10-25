# i3dwm *IN DEVELOPMENT*
Workspace helper Inspired by [https://gitlab.com/hyask/swaysome][]swaysome, making multi monitor i3 work more like dwm.

It is in development, and not yet ready for use.

Stay tuned...

## Usage

### .i3/config

Assuming you installed i3dwm with go install, and have ~/go/bin as bin path

```
# switch to workspace
bindsym $mod+1 exec /home/<user>/go/bin/i3dwm "workspace {{ ws 1 }}"
bindsym $mod+2 exec /home/<user>/go/bin/i3dwm "workspace {{ ws 2 }}"
bindsym $mod+3 exec /home/<user>/go/bin/i3dwm "workspace {{ ws 3 }}"
bindsym $mod+4 exec /home/<user>/go/bin/i3dwm "workspace {{ ws 4 }}"
bindsym $mod+5 exec /home/<user>/go/bin/i3dwm "workspace {{ ws 5 }}"
bindsym $mod+6 exec /home/<user>/go/bin/i3dwm "workspace {{ ws 6 }}"
bindsym $mod+7 exec /home/<user>/go/bin/i3dwm "workspace {{ ws 7 }}"
bindsym $mod+8 exec /home/<user>/go/bin/i3dwm "workspace {{ ws 8 }}"
bindsym $mod+9 exec /home/<user>/go/bin/i3dwm "workspace {{ ws 9 }}"
bindsym $mod+0 exec /home/<user>/go/bin/i3dwm "workspace {{ ws 10 }}"

# Move focused container to workspace
bindsym $mod+Ctrl+1 exec /home/<user>/go/bin/i3dwm "move container to workspace {{ ws 1 }}"
bindsym $mod+Ctrl+2 exec /home/<user>/go/bin/i3dwm "move container to workspace {{ ws 2 }}"
bindsym $mod+Ctrl+3 exec /home/<user>/go/bin/i3dwm "move container to workspace {{ ws 3 }}"
bindsym $mod+Ctrl+4 exec /home/<user>/go/bin/i3dwm "move container to workspace {{ ws 4 }}"
bindsym $mod+Ctrl+5 exec /home/<user>/go/bin/i3dwm "move container to workspace {{ ws 5 }}"
bindsym $mod+Ctrl+6 exec /home/<user>/go/bin/i3dwm "move container to workspace {{ ws 6 }}"
bindsym $mod+Ctrl+7 exec /home/<user>/go/bin/i3dwm "move container to workspace {{ ws 7 }}"
bindsym $mod+Ctrl+8 exec /home/<user>/go/bin/i3dwm "move container to workspace {{ ws 8 }}"
bindsym $mod+Ctrl+9 exec /home/<user>/go/bin/i3dwm "move container to workspace {{ ws 9 }}"
bindsym $mod+Ctrl+0 exec /home/<user>/go/bin/i3dwm "move container to workspace {{ ws 10 }}"
```

### Setup workspaces


i3dwm -monitor-setup <mon1> <mon2> <mon3>

Example:

```
i3dwm -monitor-setup HDMI-0 DP-5 VGA-0
```

