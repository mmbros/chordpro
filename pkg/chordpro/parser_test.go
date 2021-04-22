package chordpro

import (
	"testing"
)

const choHero = `{t:Working Class Hero}
{st:John Lennon}

[Am] [Am] 
{sot}
e|-----0--0--0--------|-------0--0--0--0---||
B|-----1--1--1--------|-------1--1--1--1---||
G|-----2--2--2--------|-------2--2--2--2---||
D|-----2--2--2--------|--0h2--2--2--2--2---||
A|-0------------------|--------------------||
E|--------------------|--------------------||
{eot}

INTRO:

[Am] [Am] [Am] [Am] 

As [Am]soon as you're [Am]born, they [G]make you feel [Am]small[Am][Am][Am]
By [Am]giving you no [Am]time, ins[G]tead of it [Am]all[Am][Am][Am]
'Til the [Am]pain is so [Am]big, you feel [G]nothing at al[Am]l[Am][Am][Am]

{soc}
[Am]Working Class [Am]Hero is [G]somethin' to be[Am][Am][Am][Am]
[Am]Working Class [G]Hero is [D]somethin' to be[Am][Am][Am][Am]
{eoc}

They hurt you at home and they hit you at school,
They hate you if you're clever and they despise a fool,
Till you're so fucking crazy you can't follow their rules,

[Am]Working Class [Am]Hero is [G]somethin' to be[Am][Am][Am][Am]
[Am]Working Class [G]Hero is [D]somethin' to be[Am][Am][Am][Am]
{comment: ( Tab from: http://www.guitartabs.cc/tabs/j/john_lennon/working_class_hero_crd_ver_3.html )}
When they've tortured and scared you for twenty odd years,
Then they expect you to pick a career,
When you can't really function you're so full of fear,

[Am]Working Class [Am]Hero is [G]somethin' to be[Am][Am][Am][Am]
[Am]Working Class [G]Hero is [D]somethin' to be[Am][Am][Am][Am]

Keep you doped with religion and sex and TV,
And you think you're so clever and classless and free,
But you're still fucking peasents as far as I can see,

[Am]Working Class [Am]Hero is [G]somethin' to be[Am][Am][Am][Am]
[Am]Working Class [G]Hero is [D]somethin' to be[Am][Am][Am][Am]

There's room at the top they are telling you still,
But first you must learn how to smile as you kill,
If you want to be like the folks on the hill,

[Am]Working Class [Am]Hero is [G]somethin' to be[Am][Am][Am][Am]
[Am]Working Class [G]Hero is [D]somethin' to be[Am][Am][Am][Am]


OUTRO:

If you [Am]wanna be a [Am]hero, well [G]just follow [Am]me

If you [Am]wanna be a [Am]hero, well [G]just follow [D]me


{sot}
e|-------------0------|-----------0--------||
B|----------1---------|----------1---------||
G|-------2------------|---------2----------||
D|----2---------------|--------2-----------||
A|-0------------------|-------0------------||
E|--------------------|--0-----------------||
{eot}

BY
SIMOSTRATO!`

func Test_ParseSimple(t *testing.T) {
	src := choHero
	ss := ParseText(src)

	t.Logf("SONGS: %s\n", ss.String())

	t.FailNow()
}
