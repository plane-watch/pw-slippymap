package bindings

import (
	"log"

	_ "golang.org/x/mobile/bind"
	_ "golang.org/x/mobile/bind/java"

	//_ "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/mobile"

	slippymap "github.com/plane-watch/pw-slippymap"
)

// Unexpported Game to avoid compile errors around multi-value returns
// Installed: ebitten mobile
//    go install github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile@latest
// Installed: android studio
// downloaded and unzipped android NDK
//    https://developer.android.com/ndk/downloads
// Embedded NDK inside SDK:
//    C:\Users\sirsquidness\appdata\Local\Android\Sdk>mklink /J ndk-bundle d:\bin\android-ndk-r23b
//    Junction created for ndk-bundle <<===>> d:\bin\android-ndk-r23b
// Set env var in VSCode PowerShell terminal:
//    $Env:ANDROID_HOME = "C:\Users\sirsquidness\AppData\Local\Android\Sdk"
//    $Env:PATH += ";d:\Program Files\Android\Android Studio\jre\bin"
// Ran build:
//    ebitenmobile.exe bind  -v  -target android -javapkg watch.plane.core -classpath watch.plane.core.gostuff -o D:\projects\planewatchapp\blah\blah.aar github.com/plane-watch/pw-slippymap/bindings
//    NOTE: d:\projects\planewatchapp\blah\blah.aar is the path to inside the android project that will exist after the first time we've imported the project
//    ebitenmobile.exe bind -x -work -v  -target android -bootclasspath watch.plane.app -classpath watch.plane.app -o D:\projects\plane-watch-app\blah.aar github.com/plane-watch/pw-slippymap/bindings
//    -x prints all of the stuff as it runs.... -work leaves the directories there afterward
//     need -javapkg watch.plane.app ?
// IMPORTANT: do not have a `vendor` directory, because ebitenmobile invokes gomobile which runs `go list -m -json all`  which doesn't work with vendored stuff
// Create new android project
// file -> new -> module -> import aar -> select the thing
// in `build.gradle` for `module: app` (not the Project gradlefile, not the newly-imported-module gradelfile) add `compile project(':blah')` to the end of dependencies
// in res/layout/activity_main.xml, open that up, and add the Layout to the View (when I did this, I had created a new class inside the project that extended the auto-generated EbitenView class from the blah.aar imported module)

func init() {
	// yourgame.Game must implement ebiten.Game interface.
	// For more details, see
	// * https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2#Game

	g, err := slippymap.GetGame()
	if err != nil {
		log.Fatalf("could not init game: %s", err)
	}
	mobile.SetGame(g)
}

// Dummy is a dummy exported function.
//
// gomobile doesn't compile a package that doesn't include any exported function.
// Dummy forces gomobile to compile this package.
func Dummy() {}
