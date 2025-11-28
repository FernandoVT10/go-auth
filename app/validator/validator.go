package validator

import "fmt"

type FieldFn func() error

func ValidateStruct(fields ...FieldFn) error {
    for _, field := range fields {
        err := field()
        if err != nil {
            return err
        }
    }

    return nil
}

type Validator func(fieldName string, value any) error

func Field(fieldName string, value any, validators ...Validator) FieldFn {
    return func() error {
        for _, validator := range validators {
            err := validator(fieldName, value)
            if err != nil {
                return err
            }
        }

        return nil
    }
}

var Required Validator = func(fieldName string, value any) error {
    switch typedValue := value.(type) {
    case string:
        if typedValue == "" {
            return fmt.Errorf("\"%s\" is required", fieldName)
        }
    default:
        fmt.Printf("[ERROR] type not implemented yet")
        return nil
    }

    return nil
}

func Length(minLen int, maxLen int) Validator {
    return func(fieldName string, value any) error {
        str, ok := value.(string)
        if !ok {
            fmt.Printf("[ERROR] Length should be used only on strings")
            return nil
        }

        if len(str) < minLen {
            return fmt.Errorf("\"%s\" min length is %d", fieldName, minLen)
        }

        if len(str) > maxLen {
            return fmt.Errorf("\"%s\" max length is %d", fieldName, maxLen)
        }

        return nil
    }
}
