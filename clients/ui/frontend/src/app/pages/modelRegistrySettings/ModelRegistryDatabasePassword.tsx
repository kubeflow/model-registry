import * as React from 'react';
import { TextField, InputAdornment, IconButton } from '@mui/material';
import { Visibility, VisibilityOff } from '@mui/icons-material';

type ModelRegistryDatabasePasswordProps = {
  password?: string;
  setPassword: (password: string) => void;
  isPasswordTouched: boolean;
  setIsPasswordTouched: (isTouched: boolean) => void;
  showPassword?: boolean;
  editRegistry?: boolean;
};

const ModelRegistryDatabasePassword: React.FC<ModelRegistryDatabasePasswordProps> = ({
  password,
  setPassword,
  isPasswordTouched,
  setIsPasswordTouched,
  showPassword,
  editRegistry,
}) => {
  const [show, setShow] = React.useState(false);

  return (
    <TextField
      required
      type={show ? 'text' : 'password'}
      value={password}
      onChange={(e) => setPassword(e.target.value)}
      onBlur={() => setIsPasswordTouched(true)}
      error={isPasswordTouched && !password?.trim().length}
      helperText={isPasswordTouched && !password?.trim().length ? 'Password cannot be empty' : ''}
      InputProps={{
        endAdornment: (
          <InputAdornment position="end">
            <IconButton onClick={() => setShow(!show)} edge="end">
              {show ? <VisibilityOff /> : <Visibility />}
            </IconButton>
          </InputAdornment>
        ),
      }}
    />
  );
};

export default ModelRegistryDatabasePassword; 