/**
 * Formats a value with optional unit for display in the hardware configuration table
 */
export const formatValue = (value: number | string | undefined, unit?: string): string => {
  if (value === undefined) {
    return '-';
  }
  if (typeof value === 'number') {
    return unit ? `${value.toFixed(1)} ${unit}` : value.toFixed(1);
  }
  return String(value);
};

/**
 * Creates a hardware configuration display string combining count and type
 */
export const formatHardwareConfiguration = (count: number, type: string): string =>
  `${count} x ${type}`;
