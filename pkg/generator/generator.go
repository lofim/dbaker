package generator

import (
	"dbaker/pkg/model"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

var (
	ErrColumnTypeNotSupported = errors.New("column type not supported")
)

/**
 * In scope: random value generation, annotation resolution, constraint solver?
 * Out of scope: database connection, writing of values.
 * Should be rather databse agnostic, might leverage a specific database writer (PostgreSQL data writer)
 */
type ValueGenerator struct{}

func (g *ValueGenerator) GenVals(cols []model.Column, iter uint32) ([]any, error) {
	// for each column generate value
	var values []any
	for _, col := range cols {
		value, err := g.GenVal(col, iter)
		if err != nil {
			return nil, fmt.Errorf("failed to generate value for column '%s(%s)': %w", col.Name, col.Typ, err)
		}

		values = append(values, value)
	}

	return values, nil
}

func (g *ValueGenerator) GenVal(col model.Column, iter uint32) (any, error) {
	if col.IsUnique {
		return g.GenUniqueVal(col, iter)
	}

	return g.GenRawVal(col)
}

// 1. look at the annotation (use annotation logic)
// 2. if no annotation try inferring meaning base on name heurestically
// 3. if no-infer tag or not possible to infer use generic type inference
// Initial implemetation does only generic type inference
func (g *ValueGenerator) GenRawVal(col model.Column) (any, error) {
	switch col.Typ {
	case model.SmallInt:
		return gofakeit.IntRange(-32768, 32767), nil
	case model.Int:
		return gofakeit.IntRange(-2147483648, 2147483647), nil
	case model.BigInt:
		return gofakeit.IntRange(-9223372036854775808, 9223372036854775807), nil
	case model.Real:
		return gofakeit.Float32(), nil
	case model.Double:
		return gofakeit.Float64(), nil

	case model.Char:
		fallthrough
	case model.Varchar:
		return gofakeit.LetterN(col.MaxLength), nil
	case model.Text:
		return gofakeit.Sentence(rand.IntN(10-0) + 1), nil

	case model.UUID:
		return gofakeit.UUID(), nil
	case model.Boolean:
		return gofakeit.Bool(), nil

	// TODO: make sure to format properly for respectice targets
	case model.Date:
		return gofakeit.Date(), nil
	case model.Time:
		return gofakeit.Date(), nil
	case model.Timestamp:
		return gofakeit.Date(), nil
	case model.TimestampTZ:
		return gofakeit.Date(), nil

	default:
		return nil, ErrColumnTypeNotSupported
	}
}

func (g *ValueGenerator) GenUniqueVal(col model.Column, iter uint32) (any, error) {
	switch col.Typ {
	case model.SmallInt:
		fallthrough
	case model.Int:
		fallthrough
	case model.BigInt:
		return iter, nil
	case model.Real:
		return float32(iter), nil
	case model.Double:
		return float64(iter), nil

	case model.Char:
		fallthrough
	case model.Varchar:
		iterDigits := uint(len(strconv.FormatUint(uint64(iter), 10)))
		return fmt.Sprintf("%d%s", iter, gofakeit.LetterN(col.MaxLength-iterDigits)), nil
	case model.Text:
		return fmt.Sprintf("%d%s", iter, gofakeit.Sentence(rand.IntN(10-1)+1)), nil

	case model.UUID:
		return gofakeit.UUID(), nil
	case model.Boolean:
		return iter%2 == 0, nil

	// TODO: make sure to format properly for respectice targets
	case model.Date:
		return time.UnixMilli(int64(iter)), nil
	case model.Time:
		return time.UnixMilli(int64(iter)), nil
	case model.Timestamp:
		return time.UnixMilli(int64(iter)), nil
	case model.TimestampTZ:
		return time.UnixMilli(int64(iter)), nil

	default:
		return nil, ErrColumnTypeNotSupported
	}
}
