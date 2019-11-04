package parser

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/oleg-balunenko/spamassassin-parser/pkg/models"
)

var (
	reType2 = regexp.MustCompile(`(?m)([-]?\d.\d)[\s]+([[:word:]]+)\s+(.*[\n]?)`)
)

type report2Parser struct{}

func (rp report2Parser) Parse(data io.Reader) (models.Report, error) {
	const (
		colFullMatch = iota
		colScore
		colTag
		colDescr
	)
	var (
		r     models.Report
		score float64
		lnum  int
		start bool
	)
	sc := bufio.NewScanner(data)
	for sc.Scan() {
		lnum++
		line := sc.Text()
		if !start {
			if strings.Contains(line, "----") {
				start = true
			}
			continue
		}

		matches := reType2.FindStringSubmatch(line)
		if len(matches) != 0 {
			sc, err := strconv.ParseFloat(matches[colScore], 64)
			if err != nil {
				return emptyReport, errors.Wrapf(err,
					"failed to parse score [line num: %d], [line: %s], score[%s]",
					lnum, line, matches[colScore])
			}

			sc = sanitizeScore(sc)
			score = score + sc
			r.SpamAssassin.Headers = append(r.SpamAssassin.Headers, models.Headers{
				Score:       sc,
				Tag:         matches[colTag],
				Description: matches[colDescr],
			})

		} else {
			last := len(r.SpamAssassin.Headers) - 1
			if last >= 0 {
				line = strings.TrimSpace(line)
				r.SpamAssassin.Headers[last].Description += " " + line
			}
		}
	}

	r.SpamAssassin.Score = sanitizeScore(score)

	return r, nil
}