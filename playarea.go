package main

import (
	"fmt"
	"github.com/powerpu/go-fake-ts"
	"github.com/spf13/cobra"
	"github.com/wcharczuk/go-chart"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var hardcodedAwesomeWebpage = `
<html>
<head>
<title>Playarea</Title>
<style>
input { width: 6em; }
input[type=button], input[type=reset] { width: 100%; }
fieldset { width: 98%; border: none; border-bottom: 2px solid black; display: inline-block; text-align: right; padding: 0; padding-top: 0.5em; margin-top: 0.5em; padding-bottom: 0.5em; margin-right: 0.3em }
legend { background-color: orange; width: 97%; padding: 0.1em; margin: 0.1em; }
</style>
</head>
<body style="text-align: center; margin: 0; padding: 0;"  onLoad="refreshChart();">
<form id="formy">
    <input id="configLine" type="text" style="width: 100%; font-size: 1.2em; float: left"/>
    <table style="border: none; padding: 0; margin: 0; width: 100%;">
        <tr>
            <td style="width: 15em; border-right: 2px solid black; vertical-align: top;">
                <fieldset>
                    <legend>General</legend>
                    <label for="startAt">Start At:</label><input type="number" id="startAt" name="startAt" value="0" onChange="refreshChart();" min="0" max="1000000"><br>
                    <label for="n">Generate Samples:</label><input type="number" id="n" name="n" value="1024" onChange="refreshChart();" min="1024" max="1000000"><br>
                    <label for="stretchStart">Stretch Start:</label><input type="number" id="stretchStart" name="stretchStart" value="1" onChange="refreshChart();" step="0.0001"><br>
                    <label for="stretchEnd">Stretch Stop:</label><input type="number" id="stretchEnd" name="stretchEnd" value="1" onChange="refreshChart();" step="0.0001"><br>
                    <label for="slope">Slope:</label><input type="number" id="slope" name="slope" value="0" onChange="refreshChart();" step="0.0001"><br>
					<label for="bump">Bump:</label><input type="number" id="bump" name="bump" value="0" onChange="refreshChart();" step="0.0001"><br>
                    <label for="from">From:</label><input type="number" id="from" name="from" value="-100" onChange="refreshChart();"><br>
                    <label for="to">To:</label><input type="number" id="to" name="to" value="100" onChange="refreshChart();"><br>
                    <label for="limitUpper">Limit Upper:</label><input type="checkbox" id="limitUpper" name="limitUpper" value="True" onClick="refreshChart();"><br>
                    <label for="limitLower">Limit Lower:</label><input type="checkbox" id="limitLower" name="limitLower" value="True" onClick="refreshChart();"><br>
                    <label for="lockRange">Lock Range:</label><input type="checkbox" id="lockRange" name="lockRange" value="True" onClick="refreshChart();"><br>
                </fieldset>
                <fieldset>
                    <legend>Permanent Bump Params</legend>
					<label for="permaBumpAt">Bump At:</label><input type="number" id="permaBumpAt" name="permaBumpAt" value="0" step="1" onChange="refreshChart();"><br>
					<label for="permaBumpBy">Bump By:</label><input type="number" id="permaBumpBy" name="permaBumpBy" value="100" step="0.1" onChange="refreshChart();"><br>
					<label for="permaBumpSmoother">Smoother:</label><input type="number" id="permaBumpSmoother" name="permaBumpSmoother" value="50" step="1" onChange="refreshChart();"><br>
                </fieldset>
                <fieldset>
                    <legend>Random Params</legend>
                    <label for="generateRandom">Generate random:</label><input type="checkbox" id="generateRandom" name="generateRandom" value="True" onClick="refreshChart();" checked><br>
					<label for="randomSeed">Seed:</label><input type="number" id="randomSeed" name="randomSeed" value="1" onChange="refreshChart();"><br>
					<label for="randomBias">Bias:</label><input type="number" id="randomBias" name="randomBias" value="0.5" onChange="refreshChart();" min="0" max="1" step="0.001"><br>
                </fieldset>
                <fieldset>
                    <legend>Spike Params</legend>
                    <label for="generateSpikes">Generate spikes:</label><input type="checkbox" id="generateSpikes" name="generateSpikes" value="True" onClick="refreshChart();"><br>
					<label for="spikeSustain">Sustain:</label><input type="number" id="spikeSustain" name="spikeSustain" value="5" onChange="refreshChart();"><br>
					<label for="spikeEvery">Spike Every:</label><input type="number" id="spikeEvery" name="spikeEvery" value="100" onChange="refreshChart();"><br>
					<label for="spikeTo">Spike To:</label><input type="number" id="spikeTo" name="spikeTo" value="100" onChange="refreshChart();"><br>
					<label for="spikeWobble">Spike Wobble:</label><input type="checkbox" id="spikeWobble" name="spikeWobble" value="True" onClick="refreshChart();"><br>
					<label for="spikeWobbleFactor">Wobble Factor:</label><input type="number" id="spikeWobbleFactor" name="spikeWobbleFactor" value="200" onChange="refreshChart();"><br>
					<label for="spikeSmoother">Smoother:</label><input type="number" id="spikeSmoother" name="spikeSmoother" value="20" onChange="refreshChart();"><br>
                </fieldset>
                <fieldset>
                    <legend>Seasonality Params</legend>
                    <label for="useSeasonality">Generate seasonality:</label> <input type="checkbox" id="useSeasonality" name="useSeasonality" value="True" onClick="refreshChart();"><br>
					<label for="seasnoalityWave1">Wave 1 Frequency:</label><input type="number" id="seasnoalityWave1" name="seasnoalityWave1" value="300" min="1" max="10000000" step="1" onClick="refreshChart();"><br>
					<label for="seasonalityWave2">Wave 2 Frequency:</label><input type="number" id="seasonalityWave2" name="seasonalityWave2" value="1" min="1" max="10000000" step="1" onClick="refreshChart();"><br>
					<label for="seasonalityWave3">Wave 3 Frequency:</label><input type="number" id="seasonalityWave3" name="seasonalityWave3" value="1" min="1" max="10000000" step="1" onClick="refreshChart();"><br>
					<label for="seasonalityWave4">Wave 4 Frequency:</label><input type="number" id="seasonalityWave4" name="seasonalityWave4" value="1" min="1" max="10000000" step="1" onClick="refreshChart();"><br>
					<label for="seasonalityWave5">Wave 5 Frequency:</label><input type="number" id="seasonalityWave5" name="seasonalityWave5" value="1" min="1" max="10000000" step="1" onClick="refreshChart();"><br>
                </fieldset>
                <fieldset>
                    <legend>Stats</legend>
                    <pre id="stats"></pre>
                </fieldset>
                <fieldset>
                    <legend>Buttons</legend>
                    <input type="button" value="Refresh" onClick="refreshChart();"><br><br>
                    <input type="button" value="Previous Set" style="width: 50%" onClick="previousSet();"><input type="button" value="Next Set" style="width: 50%" onClick="nextSet();"><br><br>
                    <input type="button" value="Randomise Seed" onClick="randomRefresh();"><br><br>
                    <input type="reset" value="Reset" onClick="refreshChart();">
                </fieldset>
            </td>
            <td style="vertical-align: top;">
                <div id="loadingMsg"></div>
                <img id="pic1" src="http:/localhost:8888/chart" style="width: 100%; border: 1px solid black"/>
            </td>
        </tr>
    </table>
</form>
<script>
function httpGetAsync(url, callback) {
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function() {
        if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
            callback(xmlHttp.responseText);
    }
    xmlHttp.open("GET", url, true); // true for asynchronous
    xmlHttp.send(null);
}

function checkPic() {
    if(!el("pic1").complete){
        setTimeout(checkPic, 1000);
    } else {
        el("pic1").style.display = '';
        el("loadingMsg").innerHTML = '';
    }
}

function updateStats(text) {
    el("stats").innerHTML = text;
}

function refreshChart() {
    var elem = el('formy').elements;
    var s = '';
    for(var i = 0; i < elem.length; i++) {
        if (elem[i].type !== "reset" && elem[i].name !== "configLine" && elem[i].name !== "") {
            if (elem[i].type === "button" || elem[i].type === "checkbox") {
                s += elem[i].name + "=" + elem[i].checked + "&";
            } else {
                s += elem[i].name + "=" + elem[i].value + "&";
            }
        }
    }
    el("pic1").src = "/chart?" + s + new Date().getTime();
    httpGetAsync("/stats?" + s + new Date().getTime(), updateStats);
    updateConfigLine();
    el("pic1").style.display = 'none';
    el("loadingMsg").innerHTML = 'Loading...';
    setTimeout(checkPic, 1000);
}

function randomRefresh() {
    el("randomSeed").value = Math.floor(Math.random() * 10000000) + 1;
    refreshChart();
}

function nextSet() {
    var a = el("startAt").value;
    var b = el("n").value;
    el("startAt").value = parseInt(a) + parseInt(b);
    refreshChart()
}

function previousSet() {
    var a = el("startAt").value;
    var b = el("n").value;
    var c = parseInt(a) - parseInt(b);
    el("startAt").value = (c < 0) ? 0 : c
    refreshChart()
}

function updateConfigLine() {
    var s = '<YOUR OWN ID>,';
    s+= el("stretchStart").value + ",";
    s+= el("stretchEnd").value + ",";
    s+= el("slope").value + ",";
    s+= el("bump").value + ",";
    s+= el("from").value + ",";
    s+= el("to").value + ",";
    s+= el("limitUpper").checked + ",";
    s+= el("limitLower").checked + ",";
    s+= el("permaBumpAt").value + ",";
    s+= el("permaBumpBy").value + ",";
    s+= el("permaBumpSmoother").value + ",";
    s+= el("generateRandom").checked + ",";
    s+= el("randomSeed").value + ",";
    s+= el("randomBias").value + ",";
    s+= el("generateSpikes").checked + ",";
    s+= el("spikeSustain").value + ",";
    s+= el("spikeEvery").value + ",";
    s+= el("spikeTo").value + ",";
    s+= el("spikeWobble").checked + ",";
    s+= el("spikeWobbleFactor").value + ",";
    s+= el("spikeSmoother").value + ",";
    s+= el("useSeasonality").checked + ",";
    s+= el("seasnoalityWave1").value + ",";
    s+= el("seasonalityWave2").value + ",";
    s+= el("seasonalityWave3").value + ",";
    s+= el("seasonalityWave4").value + ",";
    s+= el("seasonalityWave5").value + ",";
    s+= "<YOUR DESCIRPTION GOES HERE>";
    el("configLine").value = s;
}

function el(id) {
    return document.getElementById(id);
}
</script>
</body>
</html>
`

func root(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "html")
	res.Write([]byte(hardcodedAwesomeWebpage))
}

func getData(res http.ResponseWriter, req *http.Request) (*fake.Data, []float64, bool) {
	// Load one data
	startAt := getParamInt(req, "startAt", 0, 0, 1000000)
	samples := getParamInt(req, "n", 1024, 1024, 1000000)
	line := make([]string, 28)
	line[0] = "d0"
	line[1] = getParamStr(req, "stretchStart", "1")
	line[2] = getParamStr(req, "stretchEnd", "1")
	line[3] = getParamStr(req, "slope", "0")
	line[4] = getParamStr(req, "bump", "0")
	line[5] = getParamStr(req, "from", "0")
	line[6] = getParamStr(req, "to", "100")
	line[7] = getParamStr(req, "limitUpper", "FALSE")
	line[8] = getParamStr(req, "limitLower", "FALSE")
	line[9] = getParamStr(req, "permaBumpAt", "0")
	line[10] = getParamStr(req, "permaBumpBy", "100")
	line[11] = getParamStr(req, "permaBumpSmoother", "50")
	line[12] = getParamStr(req, "generateRandom", "FALSE")
	line[13] = getParamStr(req, "randomSeed", strconv.Itoa(getParamInt(req, "randomSeed", randBetween(1, 10000000), 1, 1000000)))
	line[14] = getParamStr(req, "randomBias", "2.5")
	line[15] = getParamStr(req, "generateSpikes", "FALSE")
	line[16] = getParamStr(req, "spikeSustain", "0")
	line[17] = getParamStr(req, "spikeEvery", "0")
	line[18] = getParamStr(req, "spikeTo", "0")
	line[19] = getParamStr(req, "spikeWobble", "FALSE")
	line[20] = getParamStr(req, "spikeWobbleFactor", "0")
	line[21] = getParamStr(req, "spikeSmoother", "0")
	line[22] = getParamStr(req, "useSeasonality", "FALSE")
	line[23] = getParamStr(req, "seasnoalityWave1", "50")
	line[24] = getParamStr(req, "seasonalityWave2", "40")
	line[25] = getParamStr(req, "seasonalityWave3", "30")
	line[26] = getParamStr(req, "seasonalityWave4", "20")
	line[27] = getParamStr(req, "seasonalityWave5", "10")
	lockRange := strings.EqualFold(getParamStr(req, "lockRange", "FALSE"), "true")

	data := loadOneData(line, int64(samples), true).(*fake.Data)
	vals := make([]float64, samples)
	for i := 0; i < (samples + startAt); i++ {
		data.Next()
		if i == startAt-1 {
			data.JsonStats()
		}
		if i >= startAt {
			vals[i-startAt] = data.Val().(float64)
		}
	}

	// Let's pretend this is a log message :)
	fmt.Println("---------------------------")
	fmt.Println(samples)
	fmt.Println(startAt)
	fmt.Println(line)
	fmt.Println(vals[:50])
	fmt.Println(len(vals))
	return data, vals, lockRange
}

func getStats(res http.ResponseWriter, req *http.Request) {
	data, _, _ := getData(res, req)
	res.Header().Set("Content-Type", "text/json")
	res.Write([]byte(data.JsonStats()))
}

func getChart(res http.ResponseWriter, req *http.Request) {
	data, vals, lockRange := getData(res, req)

	// Determine Y-Axis min/max values
	chartMin := data.Stats.CMin
	chartMax := data.Stats.CMax

	if lockRange {
		chartMin = data.Stats.From
		chartMax = data.Stats.To
	}

	if data.Stats.CMin < chartMin {
		chartMin = data.Stats.CMin
	}

	if chartMax < data.Stats.CMax {
		chartMax = data.Stats.CMax
	}

	// Make the X-axis,Min/Max data and time highligter data for the graph
	xVals := make([]float64, len(vals))
	minVals := make([]float64, len(vals))
	maxVals := make([]float64, len(vals))

	for i := 0; i < len(minVals); i++ {
		xVals[i] = float64(i)        // X-axis
		minVals[i] = data.Stats.From // Min data
		maxVals[i] = data.Stats.To   // Max data
	}

	// This is our main random line
	mainSeries := chart.ContinuousSeries{XValues: xVals, YValues: vals}

	// This is the slope line
	linRegSeries := &chart.LinearRegressionSeries{InnerSeries: mainSeries}

	// These are the from/to lines
	fromToLineStyle := chart.Style{Show: true, StrokeColor: chart.ColorAlternateGray}
	minSeries := &chart.ContinuousSeries{Style: fromToLineStyle, XValues: xVals, YValues: minVals}
	maxSeries := &chart.ContinuousSeries{XValues: xVals, YValues: maxVals}

	// Build the chart
	graph := chart.Chart{
		Width:  1920,
		Height: 1920,
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
			Range: &chart.ContinuousRange{
				Min: chartMin,
				Max: chartMax,
			},
		},
		Series: []chart.Series{
			mainSeries,
			linRegSeries,
			minSeries,
			maxSeries,
		},
	}

	// Serve it
	res.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, res)
}

func NewPlayareaCmd() *cobra.Command {
	var playareaCmd = &cobra.Command{
		Use:   "playarea [flags]",
		Short: "Starts a webserver playarea to quickly visualise generated data",
		Long:  `Starts a webserver which you can use to tweak parameters and see the effect they have on the random data that will be generated by the 'generate' command`,
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")
			fmt.Println("Listening on port " + fmt.Sprintf("%v", port) + "...\nOpen your browser and go to http://localhost:" + fmt.Sprintf("%v", port) + " to get a random picture.")
			http.HandleFunc("/", root)
			http.HandleFunc("/chart", getChart)
			http.HandleFunc("/stats", getStats)
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
			return nil
		},
	}

	playareaCmd.Flags().IntP("port", "p", 8080, `Port to use for playarea server.`)

	return playareaCmd
}

func getParamStr(req *http.Request, pName string, pDefaultValue string) string {
	mParam, ok := req.URL.Query()[pName]
	if ok {
		return mParam[0]
	}
	return pDefaultValue
}

func getParamInt(req *http.Request, pName string, pDefaultValue int, pMin int, pMax int) int {
	mParam, ok := req.URL.Query()[pName]
	mRet := pDefaultValue
	if ok {
		tmp, err := strconv.ParseInt(mParam[0], 10, 64) // RANDOM_SEED
		if err != nil {
			mRet = pDefaultValue
		} else if tmp < int64(pMin) || tmp > int64(pMax) {
			mRet = pDefaultValue
		} else {
			mRet = int(tmp)
		}
	}
	return mRet
}
