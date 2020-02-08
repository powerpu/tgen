# tgen
Fake time series data generator tool written in Go

## Quick documentation from within the tool itself

```
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

 'RANDOM' indicates a the literal word RANDOM on its own line indicating that 
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

For 'RANDOM' the header/columns are as follows:

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
```
