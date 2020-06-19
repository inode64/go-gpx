package gpx

import (
	"bytes"
	"encoding/xml"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	geom "github.com/twpayne/go-geom"
)

func TestWpt(t *testing.T) {
	for _, tc := range []struct {
		data          string
		wpt           *WptType
		layout        geom.Layout
		g             *geom.Point
		noTestMarshal bool
		noTestNew     bool
	}{
		{
			data: "<wpt lat=\"42.438878\" lon=\"-71.119277\"></wpt>",
			wpt: &WptType{
				Lat: 42.438878,
				Lon: -71.119277,
			},
			layout: geom.XY,
			g:      geom.NewPoint(geom.XY).MustSetCoords([]float64{-71.119277, 42.438878}),
		},
		{
			data: "<wpt lat=\"42.438878\" lon=\"-71.119277\">\n" +
				"\t<ele>44.586548</ele>\n" +
				"</wpt>",
			wpt: &WptType{
				Lat: 42.438878,
				Lon: -71.119277,
				Ele: 44.586548,
			},
			layout: geom.XYZ,
			g:      geom.NewPoint(geom.XYZ).MustSetCoords([]float64{-71.119277, 42.438878, 44.586548}),
		},
		{
			data: "<wpt lat=\"42.438878\" lon=\"-71.119277\">\n" +
				"\t<time>2001-11-28T21:05:28Z</time>\n" +
				"</wpt>",
			wpt: &WptType{
				Lat:  42.438878,
				Lon:  -71.119277,
				Time: time.Date(2001, 11, 28, 21, 5, 28, 0, time.UTC),
			},
			layout: geom.XYM,
			g:      geom.NewPoint(geom.XYM).MustSetCoords([]float64{-71.119277, 42.438878, 1006981528}),
		},
		{
			data: "<wpt lat=\"42.438878\" lon=\"-71.119277\">\n" +
				"\t<ele>44.586548</ele>\n" +
				"\t<time>2001-11-28T21:05:28Z</time>\n" +
				"\t<name>5066</name>\n" +
				"\t<desc><![CDATA[5066]]></desc>\n" +
				"\t<sym>Crossing</sym>\n" +
				"\t<type><![CDATA[Crossing]]></type>\n" +
				"</wpt>\n",
			wpt: &WptType{
				Lat:  42.438878,
				Lon:  -71.119277,
				Ele:  44.586548,
				Time: time.Date(2001, 11, 28, 21, 5, 28, 0, time.UTC),
				Name: "5066",
				Desc: "5066",
				Sym:  "Crossing",
				Type: "Crossing",
			},
			layout:        geom.XYZM,
			g:             geom.NewPoint(geom.XYZM).MustSetCoords([]float64{-71.119277, 42.438878, 44.586548, 1006981528}),
			noTestMarshal: true,
			noTestNew:     true,
		},
		{
			data: "<wpt lat=\"42.438878\" lon=\"-71.119277\">\n" +
				"\t<ele>44.586548</ele>\n" +
				"\t<time>2001-11-28T21:05:28Z</time>\n" +
				"\t<magvar>1.1</magvar>\n" +
				"\t<geoidheight>2.2</geoidheight>\n" +
				"\t<name>5066</name>\n" +
				"\t<cmt>Comment</cmt>\n" +
				"\t<desc>5066</desc>\n" +
				"\t<src>Source</src>\n" +
				"\t<link href=\"http://example.com\">\n" +
				"\t\t<text>Text</text>\n" +
				"\t\t<type>Type</type>\n" +
				"\t</link>\n" +
				"\t<sym>Crossing</sym>\n" +
				"\t<type>Crossing</type>\n" +
				"\t<fix>3d</fix>\n" +
				"\t<sat>3</sat>\n" +
				"\t<hdop>4.4</hdop>\n" +
				"\t<vdop>5.5</vdop>\n" +
				"\t<pdop>6.6</pdop>\n" +
				"\t<ageofdgpsdata>7.7</ageofdgpsdata>\n" +
				"\t<dgpsid>8</dgpsid>\n" +
				"</wpt>",
			wpt: &WptType{
				Lat:         42.438878,
				Lon:         -71.119277,
				Ele:         44.586548,
				MagVar:      1.1,
				Time:        time.Date(2001, 11, 28, 21, 5, 28, 0, time.UTC),
				GeoidHeight: 2.2,
				Name:        "5066",
				Cmt:         "Comment",
				Desc:        "5066",
				Src:         "Source",
				Link: []*LinkType{
					{
						HREF: "http://example.com",
						Text: "Text",
						Type: "Type",
					},
				},
				Sym:           "Crossing",
				Type:          "Crossing",
				Fix:           "3d",
				Sat:           3,
				HDOP:          4.4,
				VDOP:          5.5,
				PDOP:          6.6,
				AgeOfDGPSData: 7.7,
				DGPSID:        []int{8},
			},
			layout:    geom.XYZM,
			g:         geom.NewPoint(geom.XYZM).MustSetCoords([]float64{-71.119277, 42.438878, 44.586548, 1006981528}),
			noTestNew: true,
		},
	} {
		var gotWpt WptType
		assert.NoError(t, xml.Unmarshal([]byte(tc.data), &gotWpt))
		assert.Equal(t, tc.wpt, &gotWpt)
		if tc.layout != geom.NoLayout {
			assert.Equal(t, tc.g, tc.wpt.Geom(tc.layout))
		}
		if !tc.noTestMarshal {
			sb := &strings.Builder{}
			e := xml.NewEncoder(sb)
			e.Indent("", "\t")
			assert.NoError(t, e.EncodeElement(tc.wpt, xml.StartElement{Name: xml.Name{Local: "wpt"}}))
			assert.Equal(t, strings.Split(tc.data, "\n"), strings.Split(sb.String(), "\n"))
		}
		if !tc.noTestNew {
			assert.Equal(t, tc.wpt, NewWptType(tc.g))
		}
	}
}

func TestRte(t *testing.T) {
	for _, tc := range []struct {
		data          string
		rte           *RteType
		layout        geom.Layout
		g             *geom.LineString
		noTestMarshal bool
		noTestNew     bool
	}{
		{
			data: "<rte>\n" +
				"\t<rtept lat=\"42.43095\" lon=\"-71.107628\"></rtept>\n" +
				"\t<rtept lat=\"42.43124\" lon=\"-71.109236\"></rtept>\n" +
				"</rte>",
			rte: &RteType{
				RtePt: []*WptType{
					{
						Lat: 42.43095,
						Lon: -71.107628,
					},
					{
						Lat: 42.43124,
						Lon: -71.109236,
					},
				},
			},
			layout: geom.XY,
			g: geom.NewLineString(geom.XY).MustSetCoords(
				[]geom.Coord{
					{-71.107628, 42.43095},
					{-71.109236, 42.43124},
				},
			),
		},
		{
			data: "<rte>\n" +
				"\t<rtept lat=\"42.43095\" lon=\"-71.107628\">\n" +
				"\t\t<ele>23.4696</ele>\n" +
				"\t</rtept>\n" +
				"\t<rtept lat=\"42.43124\" lon=\"-71.109236\">\n" +
				"\t\t<ele>26.56189</ele>\n" +
				"\t</rtept>\n" +
				"</rte>",
			rte: &RteType{
				RtePt: []*WptType{
					{
						Lat: 42.43095,
						Lon: -71.107628,
						Ele: 23.4696,
					},
					{
						Lat: 42.43124,
						Lon: -71.109236,
						Ele: 26.56189,
					},
				},
			},
			layout: geom.XYZ,
			g: geom.NewLineString(geom.XYZ).MustSetCoords(
				[]geom.Coord{
					{-71.107628, 42.43095, 23.4696},
					{-71.109236, 42.43124, 26.56189},
				},
			),
		},
		{
			data: "<rte>\n" +
				"\t<rtept lat=\"42.43095\" lon=\"-71.107628\">\n" +
				"\t\t<time>2001-06-02T00:18:15Z</time>\n" +
				"\t</rtept>\n" +
				"\t<rtept lat=\"42.43124\" lon=\"-71.109236\">\n" +
				"\t\t<time>2001-11-07T23:53:41Z</time>\n" +
				"\t</rtept>\n" +
				"</rte>",
			rte: &RteType{
				RtePt: []*WptType{
					{
						Lat:  42.43095,
						Lon:  -71.107628,
						Time: time.Date(2001, 6, 2, 0, 18, 15, 0, time.UTC),
					},
					{
						Lat:  42.43124,
						Lon:  -71.109236,
						Time: time.Date(2001, 11, 7, 23, 53, 41, 0, time.UTC),
					},
				},
			},
			layout: geom.XYM,
			g: geom.NewLineString(geom.XYM).MustSetCoords(
				[]geom.Coord{
					{-71.107628, 42.43095, 991441095},
					{-71.109236, 42.43124, 1005177221},
				},
			),
		},
		{
			data: "<rte>\n" +
				"\t<rtept lat=\"42.43095\" lon=\"-71.107628\">\n" +
				"\t\t<ele>23.4696</ele>\n" +
				"\t\t<time>2001-06-02T00:18:15Z</time>\n" +
				"\t</rtept>\n" +
				"\t<rtept lat=\"42.43124\" lon=\"-71.109236\">\n" +
				"\t\t<ele>26.56189</ele>\n" +
				"\t\t<time>2001-11-07T23:53:41Z</time>\n" +
				"\t</rtept>\n" +
				"</rte>",
			rte: &RteType{
				RtePt: []*WptType{
					{
						Lat:  42.43095,
						Lon:  -71.107628,
						Ele:  23.4696,
						Time: time.Date(2001, 6, 2, 0, 18, 15, 0, time.UTC),
					},
					{
						Lat:  42.43124,
						Lon:  -71.109236,
						Ele:  26.56189,
						Time: time.Date(2001, 11, 7, 23, 53, 41, 0, time.UTC),
					},
				},
			},
			layout: geom.XYZM,
			g: geom.NewLineString(geom.XYZM).MustSetCoords(
				[]geom.Coord{
					{-71.107628, 42.43095, 23.4696, 991441095},
					{-71.109236, 42.43124, 26.56189, 1005177221},
				},
			),
		},
		{
			data: "<rte>\n" +
				"\t<name>BELLEVUE</name>\n" +
				"\t<desc>Bike Loop Bellevue</desc>\n" +
				"\t<number>1</number>\n" +
				"\t<rtept lat=\"42.43095\" lon=\"-71.107628\">\n" +
				"\t\t<ele>23.4696</ele>\n" +
				"\t\t<time>2001-06-02T00:18:15Z</time>\n" +
				"\t\t<name>BELLEVUE</name>\n" +
				"\t\t<cmt>BELLEVUE</cmt>\n" +
				"\t\t<desc>Bellevue Parking Lot</desc>\n" +
				"\t\t<sym>Parking Area</sym>\n" +
				"\t\t<type>Parking</type>\n" +
				"\t</rtept>\n" +
				"\t<rtept lat=\"42.43124\" lon=\"-71.109236\">\n" +
				"\t\t<ele>26.56189</ele>\n" +
				"\t\t<time>2001-11-07T23:53:41Z</time>\n" +
				"\t\t<name>GATE6</name>\n" +
				"\t\t<desc>Gate 6</desc>\n" +
				"\t\t<sym>Trailhead</sym>\n" +
				"\t\t<type>Trail Head</type>\n" +
				"\t</rtept>\n" +
				"</rte>",
			rte: &RteType{
				Name:   "BELLEVUE",
				Desc:   "Bike Loop Bellevue",
				Number: 1,
				RtePt: []*WptType{
					{
						Lat:  42.43095,
						Lon:  -71.107628,
						Ele:  23.4696,
						Time: time.Date(2001, 6, 2, 0, 18, 15, 0, time.UTC),
						Name: "BELLEVUE",
						Cmt:  "BELLEVUE",
						Desc: "Bellevue Parking Lot",
						Sym:  "Parking Area",
						Type: "Parking",
					},
					{
						Lat:  42.43124,
						Lon:  -71.109236,
						Ele:  26.56189,
						Time: time.Date(2001, 11, 7, 23, 53, 41, 0, time.UTC),
						Name: "GATE6",
						Desc: "Gate 6",
						Sym:  "Trailhead",
						Type: "Trail Head",
					},
				},
			},
			layout: geom.XYZM,
			g: geom.NewLineString(geom.XYZM).MustSetCoords(
				[]geom.Coord{
					{-71.107628, 42.43095, 23.4696, 991441095},
					{-71.109236, 42.43124, 26.56189, 1005177221},
				},
			),
			noTestNew: true,
		},
	} {
		var gotRte RteType
		assert.NoError(t, xml.Unmarshal([]byte(tc.data), &gotRte))
		assert.Equal(t, tc.rte, &gotRte)
		if tc.layout != geom.NoLayout {
			assert.Equal(t, tc.g, tc.rte.Geom(tc.layout))
		}
		if !tc.noTestMarshal {
			sb := &strings.Builder{}
			e := xml.NewEncoder(sb)
			e.Indent("", "\t")
			assert.NoError(t, e.EncodeElement(tc.rte, xml.StartElement{Name: xml.Name{Local: "rte"}}))
			assert.Equal(t, strings.Split(tc.data, "\n"), strings.Split(sb.String(), "\n"))
		}
		if !tc.noTestNew {
			assert.Equal(t, tc.rte, NewRteType(tc.g))
		}
	}
}

func TestTrk(t *testing.T) {
	for _, tc := range []struct {
		data          string
		trk           *TrkType
		layout        geom.Layout
		g             *geom.MultiLineString
		noTestMarshal bool
		noTestNew     bool
	}{
		{
			data: "<trk>\n" +
				"\t<trkseg>\n" +
				"\t\t<trkpt lat=\"47.644548\" lon=\"-122.326897\">\n" +
				"\t\t\t<ele>4.46</ele>\n" +
				"\t\t\t<time>2009-10-17T18:37:26Z</time>\n" +
				"\t\t</trkpt>\n" +
				"\t\t<trkpt lat=\"47.644548\" lon=\"-122.326897\">\n" +
				"\t\t\t<ele>4.94</ele>\n" +
				"\t\t\t<time>2009-10-17T18:37:31Z</time>\n" +
				"\t\t</trkpt>\n" +
				"\t\t<trkpt lat=\"47.644548\" lon=\"-122.326897\">\n" +
				"\t\t\t<ele>6.87</ele>\n" +
				"\t\t\t<time>2009-10-17T18:37:34Z</time>\n" +
				"\t\t</trkpt>\n" +
				"\t</trkseg>\n" +
				"</trk>",
			trk: &TrkType{
				TrkSeg: []*TrkSegType{
					{
						TrkPt: []*WptType{
							{
								Lat:  47.644548,
								Lon:  -122.326897,
								Ele:  4.46,
								Time: time.Date(2009, 10, 17, 18, 37, 26, 0, time.UTC),
							},
							{
								Lat:  47.644548,
								Lon:  -122.326897,
								Ele:  4.94,
								Time: time.Date(2009, 10, 17, 18, 37, 31, 0, time.UTC),
							},
							{
								Lat:  47.644548,
								Lon:  -122.326897,
								Ele:  6.87,
								Time: time.Date(2009, 10, 17, 18, 37, 34, 0, time.UTC),
							},
						},
					},
				},
			},
			layout: geom.XYZM,
			g: geom.NewMultiLineString(geom.XYZM).MustSetCoords(
				[][]geom.Coord{
					{
						{-122.326897, 47.644548, 4.46, 1255804646},
						{-122.326897, 47.644548, 4.94, 1255804651},
						{-122.326897, 47.644548, 6.87, 1255804654},
					},
				},
			),
		},
	} {
		var gotTrk TrkType
		assert.NoError(t, xml.Unmarshal([]byte(tc.data), &gotTrk))
		assert.Equal(t, tc.trk, &gotTrk)
		if tc.layout != geom.NoLayout {
			assert.Equal(t, tc.g, tc.trk.Geom(tc.layout))
		}
		if !tc.noTestMarshal {
			sb := &strings.Builder{}
			e := xml.NewEncoder(sb)
			e.Indent("", "\t")
			assert.NoError(t, e.EncodeElement(tc.trk, xml.StartElement{Name: xml.Name{Local: "trk"}}))
			assert.Equal(t, strings.Split(tc.data, "\n"), strings.Split(sb.String(), "\n"))
		}
		if !tc.noTestNew {
			assert.Equal(t, tc.trk, NewTrkType(tc.g))
		}
	}
}

func TestRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		data string
		gpx  *GPX
	}{
		{
			data: "<gpx" +
				" version=\"1.0\"" +
				" creator=\"ExpertGPS 1.1 - http://www.topografix.com\"" +
				" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"" +
				" xmlns=\"http://www.topografix.com/GPX/1/0\"" +
				" xsi:schemaLocation=\"http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd\">" +
				"</gpx>",
			gpx: &GPX{
				Version: "1.0",
				Creator: "ExpertGPS 1.1 - http://www.topografix.com",
			},
		},
		{
			data: "<gpx" +
				" version=\"1.0\"" +
				" creator=\"ExpertGPS 1.1 - http://www.topografix.com\"" +
				" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"" +
				" xmlns=\"http://www.topografix.com/GPX/1/0\"" +
				" xsi:schemaLocation=\"http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd\">\n" +
				"\t<wpt lat=\"42.438878\" lon=\"-71.119277\">\n" +
				"\t\t<ele>44.586548</ele>\n" +
				"\t\t<time>2001-11-28T21:05:28Z</time>\n" +
				"\t\t<name>5066</name>\n" +
				"\t\t<desc>5066</desc>\n" +
				"\t\t<sym>Crossing</sym>\n" +
				"\t\t<type>Crossing</type>\n" +
				"\t</wpt>\n" +
				"</gpx>",
			gpx: &GPX{
				Version: "1.0",
				Creator: "ExpertGPS 1.1 - http://www.topografix.com",
				Wpt: []*WptType{
					{
						Lat:  42.438878,
						Lon:  -71.119277,
						Ele:  44.586548,
						Time: time.Date(2001, 11, 28, 21, 5, 28, 0, time.UTC),
						Name: "5066",
						Desc: "5066",
						Sym:  "Crossing",
						Type: "Crossing",
					},
				},
			},
		},
		{
			data: "<gpx" +
				" version=\"1.0\"" +
				" creator=\"ExpertGPS 1.1 - http://www.topografix.com\"" +
				" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"" +
				" xmlns=\"http://www.topografix.com/GPX/1/0\"" +
				" xsi:schemaLocation=\"http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd\">\n" +
				"\t<rte>\n" +
				"\t\t<name>BELLEVUE</name>\n" +
				"\t\t<desc>Bike Loop Bellevue</desc>\n" +
				"\t\t<number>1</number>\n" +
				"\t\t<rtept lat=\"42.43095\" lon=\"-71.107628\">\n" +
				"\t\t\t<ele>23.4696</ele>\n" +
				"\t\t\t<time>2001-06-02T00:18:15Z</time>\n" +
				"\t\t\t<name>BELLEVUE</name>\n" +
				"\t\t\t<cmt>BELLEVUE</cmt>\n" +
				"\t\t\t<desc>Bellevue Parking Lot</desc>\n" +
				"\t\t\t<sym>Parking Area</sym>\n" +
				"\t\t\t<type>Parking</type>\n" +
				"\t\t</rtept>\n" +
				"\t\t<rtept lat=\"42.43124\" lon=\"-71.109236\">\n" +
				"\t\t\t<ele>26.56189</ele>\n" +
				"\t\t\t<time>2001-11-07T23:53:41Z</time>\n" +
				"\t\t\t<name>GATE6</name>\n" +
				"\t\t\t<desc>Gate 6</desc>\n" +
				"\t\t\t<sym>Trailhead</sym>\n" +
				"\t\t\t<type>Trail Head</type>\n" +
				"\t\t</rtept>\n" +
				"\t</rte>\n" +
				"</gpx>",
			gpx: &GPX{
				Version: "1.0",
				Creator: "ExpertGPS 1.1 - http://www.topografix.com",
				Rte: []*RteType{
					{
						Name:   "BELLEVUE",
						Desc:   "Bike Loop Bellevue",
						Number: 1,
						RtePt: []*WptType{
							{
								Lat:  42.43095,
								Lon:  -71.107628,
								Ele:  23.4696,
								Time: time.Date(2001, 6, 2, 0, 18, 15, 0, time.UTC),
								Name: "BELLEVUE",
								Cmt:  "BELLEVUE",
								Desc: "Bellevue Parking Lot",
								Sym:  "Parking Area",
								Type: "Parking",
							},
							{
								Lat:  42.43124,
								Lon:  -71.109236,
								Ele:  26.56189,
								Time: time.Date(2001, 11, 7, 23, 53, 41, 0, time.UTC),
								Name: "GATE6",
								Desc: "Gate 6",
								Sym:  "Trailhead",
								Type: "Trail Head",
							},
						},
					},
				},
			},
		},
	} {
		got, err := Read(bytes.NewBufferString(tc.data))
		assert.NoError(t, err)
		assert.Equal(t, tc.gpx, got)
		sb := &strings.Builder{}
		assert.NoError(t, tc.gpx.WriteIndent(sb, "", "\t"))
		assert.Equal(t, strings.Split(tc.data, "\n"), strings.Split(sb.String(), "\n"))
	}
}

func TestTime(t *testing.T) {
	for _, tc := range []struct {
		t time.Time
		m float64
	}{
		{
			t: time.Unix(0, 0),
			m: 0,
		},
		{
			t: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			m: 946684800,
		},
		{
			t: time.Date(2006, 1, 2, 15, 4, 5, 500000000, time.UTC),
			m: 1136214245.5,
		},
	} {
		assert.Equal(t, tc.m, timeToM(tc.t))
		assert.Equal(t, tc.t, mToTime(tc.m))
	}
}

func TestParseExamples(t *testing.T) {
	for _, filename := range []string{
		"testdata/ashland.gpx",
		"testdata/fells_loop.gpx",
		"testdata/mystic_basin_trail.gpx",
	} {
		t.Run(filename, func(t *testing.T) {
			//nolint:scopelint
			f, err := os.Open(filename)
			assert.NoError(t, err)
			defer f.Close()
			_, err = Read(f)
			assert.NoError(t, err)
		})
	}
}

func TestCoprightTypeYear(t *testing.T) {
	for _, tc := range []struct {
		data []byte
		year int
	}{
		{
			data: []byte("<copyright><year>2019Z</year></copyright>"),
			year: 2019,
		},
		{
			data: []byte("<copyright><year>2013</year></copyright>"),
			year: 2013,
		},
		{
			data: []byte("<copyright><year>2011+05:00</year></copyright>"),
			year: 2011,
		},
		{
			data: []byte("<copyright><year>2010-07:00</year></copyright>"),
			year: 2010,
		},
	} {
		var gotCopyright CopyrightType
		assert.NoError(t, xml.Unmarshal(tc.data, &gotCopyright))
		assert.Equal(t, tc.year, gotCopyright.Year)
	}
}
