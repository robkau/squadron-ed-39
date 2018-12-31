Squadron-ED-39
----  

This is a 2D game I created with Go. It compiles to a self contained binary executable that will run on Windows systems with no installation required. 

The game world is simulated with a simple timestepped physics and collision detection system that I created.   
[pixel](https://github.com/faiface/pixel) game library is used to render the world to an OpenGL desktop window at 60 FPS. It also provides useful vector structs and the ability to play sound.  
[Packr](https://github.com/gobuffalo/packr) is used to store music assets inside the .exe during build time.

The background music is ["Mission" by BoxCat Games](https://boxcat.bandcamp.com/track/mission) which was retrieved and shared under the [Creative Commons Attribution 3.0 Unported](https://creativecommons.org/licenses/by/3.0/) license.

---

**Story**:  
You are part of Earth Defense Squad 39 and platforms are attacking from outer space.

You were sent with one energy collector than can absorb projectiles to store 1 Joule of energy.

You can use stored energy to build turrets, they automatically fire at platforms, with some degree of accuracy.

You lose if any platforms reach the bottom of the screen.  

*Tip: Build your first turrets under the green square to guarantee steady energy supply (left click)*

---

**Controls**:  
Left Click: Place turret (costs 20 Joules)  
Right Click: Turret blast at mouse cursor (costs 15 Joules)  
r: Reset world  
F9: Profile heap and save memory dump (when in debug mode)  
*Enable debug mode by setting env var 'sq39_debug' before running client.exe, this will also save a CPU profile upon quitting*
  
---

# How to play:
####Easy option) Download binary
 - Download the game from the **releases** tab and run client.exe


####Hard option) Build the game from source
 - First you need to configure your Windows build environment for cgo, see one guide [here](https://github.com/faiface/pixel/wiki/Building-Pixel-on-Windows)  
 - Afterwards:
```
go get github.com/robkau/squadron-ed-39
cd ~/go/src/github.com/robkau/squadron-ed-39/client  
packr build  
./client.exe
```