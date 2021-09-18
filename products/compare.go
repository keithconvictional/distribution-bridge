package products

import (
	"errors"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// productsMatch custom method for comparing two products. IDs will be completely different in both.
func productsMatch(product Product, productTwo Product) error {
	// Images
	if len(product.Images) != len(productTwo.Images) {
		return errors.New("unequal number of images between both products")
	}
	foundSrcs := 0
	for _, imageFromOne := range product.Images {
		for _, imageFromTwo := range productTwo.Images {
			if imageFromOne.Src == imageFromTwo.Src {
				foundSrcs++

				if imageFromOne.Position != imageFromTwo.Position {
					return errors.New(fmt.Sprintf("image positions do not match for %s", imageFromOne.Src))
				}
				break
			}
		}
	}
	if foundSrcs != len(product.Images) {
		return errors.New(fmt.Sprintf("did not find all matches for all images (Found %d of %d)", foundSrcs, len(product.Images)))
	}

	// Variants
	if len(product.Variants) != len(productTwo.Variants) {
		return errors.New("unequal number of variants")
	}
	foundVariants := 0
	for _, variantOne := range product.Variants {
		for _, variantTwo := range product.Variants {
			if variantOne.VariantID == variantTwo.VariantID {
				foundVariants++
				if !cmp.Equal(variantOne, variantTwo, cmpopts.IgnoreFields(Variant{}, "ID")) {
					return errors.New(fmt.Sprintf("variants do not match for %s and %s", variantOne.ID, variantTwo.ID))
				}
				break
			}
		}
	}
	if foundVariants != len(product.Variants) {
		return errors.New(fmt.Sprintf("did not find all variants (Found %d of %d)", foundVariants, len(product.Variants)))
	}

	// Options
	if len(product.Options) != len(productTwo.Options) {
		return errors.New("unequal number of options")
	}
	foundOptions := 0
	for _, optionOne := range product.Options {
		for _, optionTwo := range product.Options {
			if optionOne.Name == optionTwo.Name {
				foundOptions++
				if !cmp.Equal(optionOne, optionTwo, cmpopts.IgnoreFields(Option{}, "ID")) {
					return errors.New(fmt.Sprintf("did not find matching options between %s and %s", optionOne.Name, optionTwo.Name))
				}
				break
			}
		}
	}
	if foundOptions != len(product.Options) {
		return errors.New(fmt.Sprintf("did not find all options (Found %d of %d)", foundOptions, len(product.Options)))
	}

	// All other fields that should match
	if !cmp.Equal(product, productTwo, cmpopts.IgnoreFields(Product{}, "ID", "Active", "Images", "Variants", "Options", "Created", "Updated", "CompanyObjectID","CompanyID")) {
		return errors.New("products do not match")
	}

	return nil
}