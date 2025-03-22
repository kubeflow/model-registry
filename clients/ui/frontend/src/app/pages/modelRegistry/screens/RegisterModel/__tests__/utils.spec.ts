import { RegisteredModelList } from '~/app/types';
import {
  isModelNameExisting,
  isNameValid,
} from '~/app/pages/modelRegistry/screens/RegisterModel/utils';
import { MR_CHARACTER_LIMIT } from '~/app/pages/modelRegistry/screens/RegisterModel/const';

describe('RegisterModel utils', () => {
  describe('isModelNameExisting', () => {
    const existingModelName = 'model2';
    const newModelName = 'model4';
    const modelList = {
      items: [{ name: 'model1' }, { name: existingModelName }, { name: 'model3' }],
    } as RegisteredModelList;
    it('should return true if model name exists in list', () => {
      expect(isModelNameExisting(existingModelName, modelList)).toBe(true);
    });

    it('should return false if model name does not exist in list', () => {
      expect(isModelNameExisting(newModelName, modelList)).toBe(false);
    });
  });

  describe('isNameValid', () => {
    it('should return true for valid model names (currently only limited by character count)', () => {
      expect(isNameValid('x'.repeat(MR_CHARACTER_LIMIT))).toBe(true);
      expect(isNameValid('')).toBe(true); //will be caught by form 'required' validation
    });
    it('should return false for names that are too long', () => {
      expect(isNameValid('x'.repeat(MR_CHARACTER_LIMIT + 1))).toBe(false);
    });
  });
});
