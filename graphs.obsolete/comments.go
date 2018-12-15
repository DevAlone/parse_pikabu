package graphs

import (
	. "models"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"io"
)

func GetCommentsPlotWriterTo(windowSize uint, startTimestamp int64) (io.WriterTo, error) {
	p, err := GetCommentsPlot(windowSize, startTimestamp)
	if err != nil {
		return nil, err
	}
	return p.WriterTo(10 * vg.Inch, 5 * vg.Inch, "png")
}

func SaveCommentsPlotToFile(filename string, windowSize uint, startTimestamp int64) error {
	println("drawing...")
	p, err := GetCommentsPlot(windowSize, startTimestamp)
	if err != nil {
		return err
	}
	err = p.Save(10*vg.Inch, 5*vg.Inch, filename)
	return err
}

func GetCommentsPlot(windowSize uint, startTimestamp int64) (*plot.Plot, error) {
	type Result struct {
		Bucket uint64
		X      uint64
		Y      uint32
	}
	result := []Result{}

	_, err := Db.Query(&result, `
WITH stats AS (
    SELECT MIN(creation_timestamp) as min_value, MAX(creation_timestamp) as max_value
    FROM comments WHERE creation_timestamp >= ?1
), number_of_bars_table AS (
    SELECT (stats.max_value - stats.min_value) / (?0) as number_of_bars FROM stats
)
SELECT 
    width_bucket(creation_timestamp, min_value, max_value, (CASE WHEN (number_of_bars > 0) THEN number_of_bars ELSE 1 END)) as bucket,
    MIN(creation_timestamp) as x,
    COUNT(*) AS y
FROM 
    comments, stats, number_of_bars_table
WHERE comments.creation_timestamp >= ?1
GROUP BY
    bucket
ORDER BY
    bucket
;
	`, windowSize, startTimestamp)
	if err != nil {
		return nil, err
	}

	values := make(plotter.XYs, len(result))

	for i, item := range result {
		values[i].X = float64(item.X)
		values[i].Y = float64(item.Y)
	}

	p, err := plot.New()
	if err != nil {
		return nil, err
	}
	p.Title.Text = "Комментарии"
	p.X.Label.Text = "время"
	p.Y.Label.Text = "количество"

	err = plotutil.AddLinePoints(p,
		"", values)
	if err != nil {
		return nil, err
	}

	return p, err
}
