package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"github.com/powerpu/go-fake-ts"
	"github.com/spf13/cobra"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var gDefaultConfig = `RANDOM
ID,SEED,GOOD_PCT,DESCRIPTION
r0,1,0,All bad
r1,1,0.1,10% good
r2,1,0.2,20% good
r3,1,0.3,30% good
r4,1,0.4,40% good
r5,1,0.5,50% good
r6,1,0.6,60% good
r7,1,0.7,70% good
r8,1,0.8,80% good
r9,1,0.9,90% good
r10,1,1,All good
PATTERN
ID,PATTERN_GOOD,PATTERN_BAD,DESCRIPTION
p0,0,1,All bad
p1,1,1,1 good 1 bad
p2,2,1,2 good 1 bad
p3,3,1,3 good 1 bad
p4,3,2,3 good 2 bad
p5,5,2,5 good 2 bad
p6,24,24,24 good 24 bad
p7,144,144,144 good 144 bad
p8,7,1,7 good 1 bad
p9,14,2,14 good 2 bad
p10,1,0,All good
TIMES
ID,INCREMENT,VARIANCE,DIRECTION,DESCRIPTION
t0,300000,0,0,5 minute time no variance
t1,300000,5000,0,5 minute time 5s variance
t2,300000,30000,0,5 minute time 30s variance
t3,300000,100000,0,5 minute time 1m variance
t3p,300000,0,0,5 minute time no variance
t3e,300000,100000,0,5 minute time 1m either way variance
t3a,300000,100000,1,5 minute time 1m always positive variance
t3b,300000,100000,-1,5 minute time 1m always negative variance
t4,3600000,0,0,1 hour time no variance
t5,3600000,60000,0,1 hour time 1m variance
t6,3600000,300000,0,1 hour time 5m variance
t7,3600000,600000,0,1 hour time 10m variance
t8,86400000,0,0,1 day time no variance
t9,86400000,3600000,0,1 day time 1h variance
t10,86400000,18000000,0,1 day time 5h variance
t11,86400000,36000000,0,1 day time 10h variance
DATA
ID,STRETCH_START,STRETCH_END,SLOPE,BUMP,FROM,TO,LIMIT_UPPER,LIMIT_LOWER,PERMA_BUMP_AT,PERMA_BUMP_BY,PERMA_BUMP_SMOOTHER,USE_RANDOM,RANDOM_SEED,RANDOM_BIAS,GENERATE_SPIKES,SPIKE_SUSTAIN,SPIKE_EVERY,SPIKE_TO,SPIKE_WOBBLE,SPIKE_WOBBLE_FACTOR,SPIKE_SMOOTHER,USE_SEASONALITY,SEASONALITY_WAVE1,SEASONALITY_WAVE2,SEASONALITY_WAVE3,SEASONALITY_WAVE4,SEASONALITY_WAVE5,DESCRIPTION
d0,1,1,0,0,-100,100,false,false,0,0,0,true,1,0.5,false,5,100,100,false,200,20,false,300,1,1,1,1,Random example
d1,1,1,0,0,-100,100,false,false,43,100,1,true,1,0.5,false,5,100,100,false,200,20,false,300,1,1,1,1,Random example with a permanent bump
d2,1,1,0,0,-100,100,false,false,0,0,0,false,1,0.5,false,5,100,100,false,200,20,true,300,1,1,1,1,Seasonality example
d3,1,1,0,0,-100,100,false,false,0,0,0,true,1,0.5,true,5,200,100,false,200,20,false,300,1,1,1,1,Random example with a spike
d4,1,1,0,100,-100,100,false,false,0,0,0,true,1,0.5,true,5,200,-100,false,200,20,false,300,1,1,1,1,Random example with a negative spike
d5,1,100,0,0,-100,100,false,false,0,0,0,false,1,0.5,false,5,100,100,false,200,20,true,300,1,1,1,1,Seasonality example with a stretch
d6,100,1,0,0,-100,100,false,false,0,0,0,false,1,0.5,false,5,100,100,false,200,20,true,300,1,1,1,1,Seasonality example with a squish
`

var gDefaultTemplate = `example,uuid={{ uuid }} [r5 random]({{ if $.r5.Val }}TRUE{{ else }}FALSE{{ end }}) [p2 Pattern]({{ $.p2.Val }}) [d0 Value]({{$.d0.Val}}) [d1 Value]({{$.d1.Val}}) [d2 Value]({{$.d2.Val}}) [d3 Value]({{$.d3.Val}}) [t3b Time Epoch]({{$.t3b.Val.Unix}})  [t3b Time]({{$.t3b.Val}})
`

func NewGenerateCmd() *cobra.Command {
	var generateCmd = &cobra.Command{
		Use:   "generate [flags]",
		Short: "Generates fake time series data",
		Long:  `Generates fake time series data based on configuration and templates`,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, _ := cmd.Flags().GetString("config")
			samples, _ := cmd.Flags().GetInt("samples")
			stream, _ := cmd.Flags().GetBool("stream")
			tmplt, _ := cmd.Flags().GetString("template")
			fromTime, _ := cmd.Flags().GetString("fromTime")
			fromTimeFormat, _ := cmd.Flags().GetString("fromTimeFormat")
			fromTimeZone, _ := cmd.Flags().GetString("fromTimeZone")
			offset, _ := cmd.Flags().GetInt("offset")
			out, _ := cmd.Flags().GetString("out")
			rate, _ := cmd.Flags().GetInt("rate")
			jitter, _ := cmd.Flags().GetInt("jitter")
			stats, _ := cmd.Flags().GetInt("stats")
			keepStats := stats > 0

			if stream {
				samples = math.MaxInt64
			}

			parsedConfig := loadConfig(config, math.MaxInt64, fromTime, fromTimeFormat, fromTimeZone, keepStats)

			// Do the actual printing of data
			t, err := template.New("default").Parse("")
			check(err)
			t = t.Funcs(template.FuncMap{"array": convertToArray})
			t = t.Funcs(template.FuncMap{"toNano": convertToNano})
			t = t.Funcs(template.FuncMap{"toSeconds": convertToSeconds})
			// t = t.Funcs(template.FuncMap{"progress": printProgress})
			t = t.Funcs(template.FuncMap{"floatToInt": convertToInt})
			t = t.Funcs(template.FuncMap{"uuid": generateUUID})
			t = t.Funcs(template.FuncMap{"seq": generateSequence})
			_, err = t.Parse(gDefaultTemplate)
			check(err)

			if tmplt != "" {
				t, err = t.ParseFiles(tmplt)
			}
			check(err)
			mTmp := strings.Split(tmplt, "/")

			var t2 *template.Template
			if tmplt != "" {
				t2 = t.Lookup(mTmp[len(mTmp)-1])
			} else {
				t2 = t.Lookup("default")
			}

			outW := bufio.NewWriter(os.Stdout)
			if out != "-" {
				f, err := os.OpenFile(out, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
				defer f.Close()
				check(err)
				sleep(rand.Float64(), rate, jitter)
				outW = bufio.NewWriter(f)
			}

			i := int64(0)
			for ; i < int64(samples+offset); i++ {
				if keepStats && i != 0 && i%int64(stats) == 0 {
					printStats(i, int64(stats), outW, parsedConfig)
				}

				for _, v := range parsedConfig {
					v.Next()
				}

				if keepStats {
					continue
				}

				if i >= int64(offset) {
					t2.Execute(outW, parsedConfig)
					sleep(rand.Float64(), rate, jitter)
					outW.Flush()
				}
			}

			// Capture anything missed in the loop
			if keepStats {
				printStats(i, int64(stats), outW, parsedConfig)
			}

			return nil
		},
	}

	generateCmd.Flags().StringP("config", "c", "", `Which config file to use?

You can generate a sample config and templates using the 'config' command.

The config file is a simple CSV file in the following format:

RANDOM
[RANDOM HEADER]
{RANDOM ENTRIES}
PATTERN
[PATTERN HEADER]
{PATTERN ENTRIES}
TIMES
[TIMES HEADER]
{TIMES ENTRIES}
DATA
[DATA HEADER]
{DATA ENTRIES}

Where:

 'RANDOM' indicates a the literal word RANOM on its own line indicating that 
         the lines that follow are loading random entries.

 'PATTERN' indicates a the literal word TIMES on its own line indicating that 
         the lines that follow are loading pattern entries.

 'TIMES' indicates a the literal word TIMES on its own line indicating that the
         lines that follow are loading time entries.

 'DATA'  indicates a the literal word DATA on its own line indicating that the
         lines that follow are loading data entries.


RANDOM
======

An example RANDOM entry may look like the following:

 r3,1,0.3,30% good

For 'RANOM' the header/columns are as follows:

 ID
    An id for this entry so that you can refer to it in the template like this
    '$.r0.Val'

 SEED
    A seed to use for any random numbers used internally

 GOOD_PCT
    A percentage of samples that will return 'true'

 DESCRIPTION
    Optional description to help the user


PATTERN
=======

An example PATTERN entry may look like the following:

 p2,2,1,2 good 1 bad

For 'PATTERN' the header/columns are as follows:

 ID
    An id for this entry so that you can refer to it in the template like this
    '$.p0.Val'

 PATTERN_GOOD
    Starting at 0, how many good samples to generate

 PATTERN_BAD
    After the good samples, how many "bad" sampled to generate.  For example if
    you want every third sample to be "bad" you would set PATTERN_GOOD to 2 and
    PATTERN_BAD to 1. You can then use it in the template like this

          '{{- if $.p0.Val -}}'

    to determine whether this is a 'true' value.

 DESCRIPTION
    Optional description to help the user


TIMES
=====

An example TIME entry may look like the following:

 t3e,300000,100000,0,5 minute time 1m either way variance

For 'TIMES' the header/columns are as follows:

 ID
    An id for this entry so that you can refer to it in the template like this
    '$.t3e.Val'

 INCREMENT
    Increment (in milliseconds) for the next entry

 VARIANCE
	When the next value is calculated, you can use a variance to add/subtract
	up to the amount of milliseconds to that time. If you want perfect times
	set this to 0.

 DIRECTION
	When the next value is calculated, you can use a value less than 0 to
	indicate the variance will always subtract from the intended time, a value
	of 0 to either subtract or add at a 50% chance or a value above 0 to always
	add to the intended time.

	In other words if time is 12:00:00PM, variance is 5 seconds and direction
	-1 this could yield 11:59:57, but never more than 12:00:00PM. Similarly, if
	direction is 1 it could yield 12:00:03 but never less than 12:00:00PM. And
	finally, if direction is 0 then either case is equally probable.

 DESCRIPTION
    Optional description to help the user

DATA
====

An example DATA entry may look like the following:

 d2,1,1,0,0,-100,100,false,false,43,100,1,true,1,0.5,false,5,100,100,false,200,20,false,300,1,1,1,1,0 to 100 soft lmit upward slope likely to breach max

For 'DATA' the header/columns are as follows:

 ID
    An id for this entry so that you can refer to it in the template like this
    '$.d2.Val'

 STRETCH_START
 STRETCH_END
	When a default graph is generated this value "stretches" or "squishes" the
	data up and down. For example if the generated data has a minimum of 35 and
	a maximum of 60, "stretching" it will increase the difference between the
	minumum and maximum while "squishing" will decrease it.

	A value above 1 will stretch, a value below 1 but more than 0 will "squish"
	while a value of 1 will be "normal". 

	STRETCH_START and STRETCH_END values indicates the start and stop values.
	This is useful if you want to "amplify over time" or "dampen over time".

	For example if STRETCH_START is 100, STRETCH_END is 0 and we want 100
	samples then we will gradually reduce the stretch from 100 to 0 over 100
	linear steps (i.e. the stretch value at the 50th sample will be 50).

 SLOPE
	This will set the slope of the generated data. 0 means whatever the
	underlying data is. Positive values means data will trend upwards, negative
	values mean data will trend downward.

 BUMP
	Tweaking this number will "bump" the random value up (or down if negative).
	For example if a maximum is 60, a positive "bump" will increase this value,
	a negative "bump" will decrease it. This is useful if you want to ensure
	values breach min/max in certain cases.

 FROM
    Generate data "from" this number. E.g. CPU values may be from 0 to 100.

 TO
    Generate data "to" this number. E.g. CPU values may be from 0 to 100.

 LIMIT_UPPER
    If numbers go above "to", when this is "TRUE" values will be set to the
    "to" value. E.g. a CPU cannot go above 100%.

 LIMIT_LOWER
    If numbers go below "from", when this is "TRUE" values will be set to the
    "from" value. E.g. a CPU cannot go below 0%.

 PERMA_BUMP_AT
	Generate data "from" this number. E.g. CPU values may be from 0 to 100. Use
	0 to disable permanent bump.

 PERMA_BUMP_BY
	What value should we bump to, expressed as a percentage of "to". Can be
	negative too.

 PERMA_BUMP_SMOOTHER
	When smoother is 1 we will go straight from current value to bump value.
	When smoother is above 1 when we will take this number of samples to reach
	the bump value (i.e. gradually but quickly rise to a spike). Use 0 to
	disable permanent bump.

 USE_RANDOM
    Whether to generate random numbers. When "TRUE" numbers will be generated
    based on the below parameters.

 RANDOM_SEED
    A seed to use for any random numbers used internally

 RANDOM_BIAS
    A parameter between 0 and 1. Use it to control the range and slope of your
    data. Generally lower biases make the slope negative but which values
	entirely depend on the underlying dataset. As always, have a play to
    see the effects.

 GENERATE_SPIKES
    Whether to generate spikes. When "TRUE" spikes will be generated based on
    the below parameters.

 SPIKE_SUSTAIN
    When a spike is reached, for how many samples should we sustain it.

 SPIKE_EVERY
    Starting at 0, every n samples will reach the desired spike value.

 SPIKE_TO
    What value should we spike to, expressed as a percentage of "to".

 SPIKE_WOBBLE
    When we're spiking and sustaining it, do we use a flat value at the top or
    do we "bounce off the top" a little for a little variation? When "TRUE" it
	means that when sustaining we'll add a little variance. We respect and will
	never break the SPIKE_TO value though!

 SPIKE_WOBBLE_FACTOR
    A magic number to tweak the "wobbliness". Have a play with this value.
	Generally a higher value will mean smoother values.

 SPIKE_SMOOTHER
    When smoother is 0 we will go straight from current value to spike value.
    When smoother is above 0 when we will take this number of samples to reach
    the spike value (i.e. gradually but quickly rise to a spike)

 USE_SEASONALITY
    Whether to generate seasonality using SIN. When "TRUE" waves will be
    generated using the WAVE parameters below.

 SEASONALITY_WAVE[1-5]
    Indicates number of points where one SIN cycle will be complete.  Each wave
    is summed to generate interference.

 DESCRIPTION
    Optional, but highly desirable, description to help the user.

You can use the 'playarea' command to play with the values and generate the
corresponding DATA entries.

Any row starting with a '#' will be treated as a comment.`)

	generateCmd.Flags().IntP("samples", "n", 10, `Number of samples to generate.`)

	generateCmd.Flags().BoolP("stream", "s", false, `Continuously generate data.`)

	generateCmd.Flags().StringP("template", "t", "", `Template to use.
Uses the Go template language (See https://golang.org/pkg/text/template/)`)

	generateCmd.Flags().StringP("fromTime", "f", "now()", `Date and/or time to use as a starting point for *all* time generation.
Uses Go date format, see https://golang.org/pkg/time/#Parse.`)

	generateCmd.Flags().StringP("fromTimeFormat", "d", "Mon Jan 2 15:04:05 MST 2006", `Date format to use for the 'from' parameter.
The default format is the output of the 'date' command.`)

	generateCmd.Flags().StringP("fromTimeZone", "z", "", `Timezone to use.
See https://en.wikipedia.org/wiki/List_of_tz_database_time_zones.`)

	generateCmd.Flags().StringP("out", "p", "-", `Print to file. Use '-' for STDOUT.`)

	generateCmd.Flags().IntP("offset", "o", 0, `Generates data from a particular starting point. Used for continuing runs.`)

	generateCmd.Flags().IntP("rate", "r", 0, `What delay in milliseconds to use between each sample.
This is useful if you need to run this continuously but slowly as opposed to
one large dump.`)

	generateCmd.Flags().IntP("jitter", "j", 0, `What "jitter" to add to the rate in milliseconds.
At random a number between 0 and value submitted will be added or subtracted at
random for a variance in rate of output. For example, if rate is 100 and jitter
is 20 actual rate will be between 80 and 120.`)

	generateCmd.Flags().IntP("stats", "m", 0, `Print the stats of the data in
the input config not the numbers. This is useful if you want to get a feel
for what the data will look like.`)

	return generateCmd
}

func printStats(i, stats int64, outW *bufio.Writer, parsedConfig map[string]fake.Value) {
	for _, v := range parsedConfig {
		outW.Write([]byte(v.JsonStats() + "\n"))
		outW.Flush()
	}
}

func sleep(n float64, rate, jitter int) {
	if rate == 0 {
		return
	}

	var tmp int

	if n <= 0.5 {
		tmp = rate + rand.Intn(jitter+1)
	} else {
		tmp = rate - rand.Intn(jitter+1)
	}

	if tmp < 0 {
		tmp = 0
	}

	time.Sleep(time.Duration(int32(tmp)) * time.Millisecond)
}

func loadConfig(config string, samples int64, fromTime string, fromTimeFormat string, fromTimeZone string, keepStats bool) map[string]fake.Value {

	const fRandom int = 0
	const fPattern int = 1
	const fTimes int = 2
	const fData int = 3

	var mReader *csv.Reader
	if config == "" {
		mReader = csv.NewReader(strings.NewReader(gDefaultConfig))
	} else {
		mCSVFile, err := os.Open(config)
		check(err)
		mReader = csv.NewReader(bufio.NewReader(mCSVFile))
	}

	mReader.FieldsPerRecord = -1 // Don't raise a "wrong number of fields" error
	mode := -1
	out := make(map[string]fake.Value)

	for {
		line, err := mReader.Read()
		if err != io.EOF {
			check(err)
		} else {
			break
		}

		// check if first character is a "#". We treat those as comments. If
		// the first word is 'ID' then it's the header row and we skip that
		// too.
		if strings.HasPrefix(line[0], "#") || strings.HasPrefix(line[0], "ID") {
			continue
		}

		// check what mode we're currently parsing
		switch line[0] {
		case "RANDOM":
			mode = fRandom
			continue
		case "PATTERN":
			mode = fPattern
			continue
		case "TIMES":
			mode = fTimes
			continue
		case "DATA":
			mode = fData
			continue
		}

		switch mode {
		case fRandom:
			out[line[0]] = loadOneRandom(line, keepStats)
		case fPattern:
			out[line[0]] = loadOnePattern(line, keepStats)
		case fTimes:
			out[line[0]] = loadOneTime(line, fromTime, fromTimeFormat, fromTimeZone, keepStats)
		case fData:
			out[line[0]] = loadOneData(line, samples, keepStats)
		}
	}

	return out
}

func loadOneRandom(line []string, keepStats bool) fake.Value {
	id := line[0]                                  // id
	seed, err := strconv.ParseInt(line[1], 10, 64) // seed
	check(err)
	pctGood, err := strconv.ParseFloat(line[2], 64) // pctGood
	check(err)
	fr, err := fake.NewRandom(id, seed, pctGood, keepStats)
	check(err)
	return fr
}

func loadOnePattern(line []string, keepStats bool) fake.Value {
	id := line[0]                                         // id
	patternGood, err := strconv.ParseInt(line[1], 10, 32) // patternGood
	check(err)
	patternBad, err := strconv.ParseInt(line[2], 10, 32) // patternBad
	check(err)
	fp, err := fake.NewPattern(id, int32(patternGood), int32(patternBad), keepStats)
	check(err)
	return fp
}

func loadOneTime(line []string, fromTime string, fromTimeFormat string, fromTimeZone string, keepStats bool) fake.Value {

	id := line[0]                           // id
	increment, err := strconv.Atoi(line[1]) // increment
	check(err)
	variance, err := strconv.Atoi(line[2]) // variance
	check(err)
	direction, err := strconv.Atoi(line[3]) // direction
	check(err)

	// In case we have a default date format but set the timezone then we'll
	// tweak the date format to remove the timezone provided
	if fromTimeZone != "" && fromTimeFormat == "Mon Jan 2 15:04:05 MST 2006" {
		fromTimeFormat = "Mon Jan 2 15:04:05 2006"
	}

	if fromTime == "now()" {
		fromTime = time.Now().Format(fromTimeFormat)
	}

	from, err := time.Parse(fromTimeFormat, fromTime)
	check(err)
	if fromTimeZone != "" {
		loc, err := time.LoadLocation(fromTimeZone)
		check(err)
		from, err = time.ParseInLocation(fromTimeFormat, fromTime, loc)
		check(err)
	}

	ft, err := fake.NewTime(id, from, increment, variance, direction, keepStats)
	check(err)
	return ft
}

func loadOneData(line []string, samples int64, keepStats bool) fake.Value {
	id := line[0] // id

	stretchStart, err := strconv.ParseFloat(line[1], 64) // STRETCH_START
	check(err)

	stretchEnd, err := strconv.ParseFloat(line[2], 64) // STRETCH_END
	check(err)

	slope, err := strconv.ParseFloat(line[3], 64) // SLOPE
	check(err)

	bump, err := strconv.ParseFloat(line[4], 64) // BUMP
	check(err)

	from, err := strconv.ParseFloat(line[5], 64) // FROM
	check(err)

	to, err := strconv.ParseFloat(line[6], 64) // TO
	check(err)

	limitUpper, err := strconv.ParseBool(line[7]) // LIMIT_UPPER
	check(err)

	limitLower, err := strconv.ParseBool(line[8]) // LIMIT_LOWER
	check(err)

	permaBumpAt, err := strconv.ParseInt(line[9], 10, 64) // PERMA_BUMP_AT
	check(err)

	permaBumpBy, err := strconv.ParseFloat(line[10], 64) // PERMA_BUMP_BY
	check(err)

	permaBumpSmoother, err := strconv.ParseInt(line[11], 10, 64) // PERMA_BUMP_SMOOTHER
	check(err)

	useRandom, err := strconv.ParseBool(line[12]) // USE_RANDOM
	check(err)

	seed, err := strconv.ParseInt(line[13], 10, 64) // RANDOM_SEED
	check(err)

	bias, err := strconv.ParseFloat(line[14], 64) // RANDOM_BIAS
	check(err)

	spike, err := strconv.ParseBool(line[15]) // GENERATE_SPIKES
	check(err)

	spikeSustain, err := strconv.ParseInt(line[16], 10, 64) // SPIKE_SUSTAIN
	check(err)

	spikeEvery, err := strconv.ParseInt(line[17], 10, 64) // SPIKE_EVERY
	check(err)

	spikeTo, err := strconv.ParseInt(line[18], 10, 64) // SPIKE_TO
	check(err)

	spikeWobble, err := strconv.ParseBool(line[19]) // SPIKE_WOBBLE
	check(err)

	spikeWobbleFactor, err := strconv.ParseInt(line[20], 10, 64) // SPIKE_WOBBLE_FACTOR
	check(err)

	spikeSmoother, err := strconv.ParseInt(line[21], 10, 64) // SPIKE_SMOOTHER
	check(err)

	seasonality, err := strconv.ParseBool(line[22]) // USE_SEASONALITY
	check(err)

	seasonalityWave1, err := strconv.ParseInt(line[23], 10, 64) // SEASONALITY_WAVE1
	check(err)

	seasonalityWave2, err := strconv.ParseInt(line[24], 10, 64) // SEASONALITY_WAVE2
	check(err)

	seasonalityWave3, err := strconv.ParseInt(line[25], 10, 64) // SEASONALITY_WAVE3
	check(err)

	seasonalityWave4, err := strconv.ParseInt(line[26], 10, 64) // SEASONALITY_WAVE4
	check(err)

	seasonalityWave5, err := strconv.ParseInt(line[27], 10, 64) // SEASONALITY_WAVE5
	check(err)

	fd, err := fake.NewData(
		id,
		samples,

		stretchStart,
		stretchEnd,
		slope,
		bump,
		from,
		to,
		limitUpper,
		limitLower,

		permaBumpAt,
		permaBumpBy,
		permaBumpSmoother,

		useRandom,
		seed,
		bias,

		spike,
		spikeEvery,
		spikeSustain,
		spikeTo,
		spikeWobble,
		spikeWobbleFactor,
		spikeSmoother,

		seasonality,
		seasonalityWave1,
		seasonalityWave2,
		seasonalityWave3,
		seasonalityWave4,
		seasonalityWave5,

		keepStats)
	check(err)
	return fd
}

func convertToNano(args ...interface{}) int64 {
	if len(args) == 1 {
		return args[0].(time.Time).UnixNano()
	}

	if rand.Float64() < 0.5 {
		return args[1].(time.Time).UnixNano() - int64(randBetween(0, args[0].(int)))
	}

	return args[1].(time.Time).UnixNano() + int64(randBetween(0, args[0].(int)))
}

func convertToSeconds(args ...interface{}) int64 {
	if len(args) == 1 {
		return args[0].(time.Time).Unix()
	}

	if rand.Float64() < 0.5 {
		return args[1].(time.Time).Unix() - int64(randBetween(0, args[0].(int)))
	}

	return args[1].(time.Time).Unix() + int64(randBetween(0, args[0].(int)))
}

// func printProgress(args ...interface{}) string {
// 	mCurrent := args[0].(int)
// 	if mCurrent%10000 == 0 && mCurrent > 0 {
// 		mLen := len(args[1].([]time.Time))
// 		os.Stderr.WriteString("Printed " + strconv.Itoa(mCurrent) + " lines out of " + strconv.Itoa(mLen) + " (" + (strconv.FormatFloat((float64(mCurrent)/float64(mLen))*100, 'f', 2, 64)) + "%).\n")
// 	}
// 	return ""
// }

func convertToArray(args ...interface{}) []interface{} {
	return args
}

func convertToInt(args ...interface{}) int64 {
	return int64(args[0].(float64))
}

func generateUUID(args ...interface{}) string {
	return uuid.New().String()
}

func generateSequence(args ...interface{}) []interface{} {
	ints := true
	switch args[0].(type) {
	case float32:
		ints = false
	case float64:
		ints = false
	}

	out := make([]interface{}, 0)
	if ints {
		from, err := strconv.ParseInt(fmt.Sprintf("%v", args[0]), 10, 64)
		if err == nil {
			to, err1 := strconv.ParseInt(fmt.Sprintf("%v", args[1]), 10, 64)
			step, err2 := strconv.ParseInt(fmt.Sprintf("%v", args[2]), 10, 64)

			if err1 != nil || err2 != nil {
				return out
			}

			for i := from; i <= to; i = i + step {
				out = append(out, i)
			}

		}
	} else {
		from, err := strconv.ParseFloat(fmt.Sprintf("%v", args[0]), 64)
		if err == nil {
			to, err1 := strconv.ParseFloat(fmt.Sprintf("%v", args[1]), 64)
			step, err2 := strconv.ParseFloat(fmt.Sprintf("%v", args[2]), 64)

			if err1 != nil || err2 != nil {
				return out
			}

			for i := from; i <= to; i = i + step {
				out = append(out, i)
			}
		}
	}

	return out
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		panic(e)
		os.Exit(1)
	}
}

func randBetween(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
