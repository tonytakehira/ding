/*
Copyright © 2021 Tony Takehira <tony@竹平.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"log"
	"time"
	"image/color"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var cfgFile string
var wantsLouder bool

const (
	screenWidth  = 128 + 64
	screenHeight = 64
	sampleRate   = 44100
)

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

type Game struct {
	audioContext *audio.Context
	audioPlayer  *audio.Player
}

var g Game
var done = make(chan struct{})
var audioDone = make(chan struct{})
var endtimer = make(chan int, 2)


// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ding",
	Short: "A command line timer",
	Long: `By default this timer uses seconds before playing a small ding sound.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		go func() {
			select {
			case <-audioDone:
				for {
					if (g.audioPlayer.IsPlaying()) {
						continue
					} else {
						os.Exit(1)
					}
				}
			}
		}()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ding.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&wantsLouder, "loud", "l", false, "use a louder alarm sound")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("loud", "l", false, "use a louder alarm sound")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ding" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ding")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

}

func gameInit() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	// Initialize audio context.
	g.audioContext = audio.NewContext(sampleRate)

	// In this example, embedded resource "Jab_wav" is used.
	//
	// If you want to use a wav file, open this and pass the file stream to wav.Decode.
	// Note that file's Close() should not be closed here
	// since audio.Player manages stream state.
	//
	f, err := os.Open("audio/hand-bell.wav")
	if wantsLouder {
		f, err = os.Open("audio/alarm.wav")
	}
	if err != nil {
		log.Fatal(err)
	}

	d, err := wav.Decode(g.audioContext, f)
	//     ...

	// Decode wav-formatted data and retrieve decoded PCM stream.
	// d, err := wav.Decode(g.audioContext, bytes.NewReader(raudio.Jab_wav))
	if err != nil {
		log.Fatal(err)
	}

	// Create an audio.Player that has one stream.
	g.audioPlayer, err = audio.NewPlayer(g.audioContext, d)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	select {
	case <- done:
		g.audioPlayer.Rewind()
		g.audioPlayer.Play()
		close(audioDone)
	default:
		return nil
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	e := <-endtimer
	HOURS := 60*60
	MINUTES := 60
	hours := e / (HOURS)
	minutes := e % (HOURS) / MINUTES
	seconds := e % MINUTES
	if g.audioPlayer.IsPlaying() {
		ebitenutil.DebugPrint(screen, "Bump!")
	} else {
		const x, y = 20, 40
		text.Draw(screen, fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds), mplusBigFont, x, y, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func openGame() {
	gameInit()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ding")
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

func startCountDown(seconds int) {
	waitTime := seconds
	endtimer <- waitTime
	for {
		select {
		case <-time.After(1 * time.Second):
			endtimer <- waitTime - 1
			waitTime--
			if waitTime <= 0 {
				close(done)
			}
		case <-done:
			break
		}

	}
}

func waitForExit() {
	select {
	case <-audioDone:
		for {
			if (g.audioPlayer.IsPlaying()) {
				continue
			} else {
				os.Exit(1)
			}
		}
	}
}