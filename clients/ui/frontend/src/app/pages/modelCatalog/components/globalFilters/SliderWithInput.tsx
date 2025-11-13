import * as React from 'react';
import {
  Flex,
  FlexItem,
  InputGroup,
  InputGroupItem,
  InputGroupText,
  Slider,
  TextInput,
} from '@patternfly/react-core';

type SliderWithInputProps = {
  value: number;
  min: number;
  max: number;
  isDisabled: boolean;
  onChange: (value: number) => void;
  suffix: string;
  ariaLabel: string;
  shouldRound?: boolean;
  showBoundaries?: boolean;
  hasTooltipOverThumb?: boolean;
};

const SliderWithInput: React.FC<SliderWithInputProps> = ({
  value,
  min,
  max,
  isDisabled,
  onChange,
  suffix,
  ariaLabel,
  shouldRound = false,
  showBoundaries = false,
  hasTooltipOverThumb = false,
}) => {
  const roundValue = React.useCallback(
    (val: number) => (shouldRound ? Math.round(val) : val),
    [shouldRound],
  );

  const [inputValue, setInputValue] = React.useState<string>(() => String(roundValue(value)));

  React.useEffect(() => {
    setInputValue(String(roundValue(value)));
  }, [value, roundValue]);

  const handleInputChange = (_event: React.FormEvent<HTMLInputElement>, val: string) => {
    setInputValue(val);

    if (val !== '') {
      const numValue = Number(val);
      if (!Number.isNaN(numValue)) {
        const processedValue = roundValue(numValue);
        const isInRange = processedValue >= min && processedValue <= max;
        if (isInRange) {
          onChange(processedValue);
        }
      }
    }
  };

  const handleSliderChange = (_event: unknown, val: number) => {
    onChange(roundValue(val));
  };

  const handleBlur = () => {
    const numValue = inputValue === '' ? min : Number(inputValue);
    const clampedValue = Number.isNaN(numValue) ? min : Math.min(Math.max(numValue, min), max);
    const rounded = roundValue(clampedValue);

    onChange(rounded);
    setInputValue(String(rounded));
  };

  return (
    <Flex alignItems={{ default: 'alignItemsCenter' }} spaceItems={{ default: 'spaceItemsMd' }}>
      <FlexItem flex={{ default: 'flex_1' }} style={{ minWidth: '300px' }}>
        <Slider
          min={min}
          max={max}
          value={value}
          onChange={handleSliderChange}
          isInputVisible={false}
          isDisabled={isDisabled}
          showBoundaries={showBoundaries}
          hasTooltipOverThumb={hasTooltipOverThumb}
        />
      </FlexItem>
      <FlexItem>
        <InputGroup>
          <InputGroupItem isFill>
            <TextInput
              type="number"
              value={inputValue}
              min={min}
              max={max}
              isDisabled={isDisabled}
              onChange={handleInputChange}
              onBlur={handleBlur}
              style={{ width: '80px' }}
              aria-label={ariaLabel}
            />
          </InputGroupItem>
          <InputGroupText isDisabled={isDisabled}>{suffix}</InputGroupText>
        </InputGroup>
      </FlexItem>
    </Flex>
  );
};

export default SliderWithInput;
