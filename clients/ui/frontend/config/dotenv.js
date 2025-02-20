const fs = require('fs');
const path = require('path');
const dotenv = require('dotenv');
const dotenvExpand = require('dotenv-expand');
const Dotenv = require('dotenv-webpack');

/**
 * Determine if the project is standalone or nested.
 *
 * @param {string} directory
 * @returns {boolean}
 */
const getProjectIsRootDir = (directory) => {
  const dotenvLocalFile = path.resolve(directory, '.env.local');
  const dotenvFile = path.resolve(directory, '.env');
  let localIsRoot;
  let isRoot;

  if (fs.existsSync(dotenvLocalFile)) {
    const { IS_PROJECT_ROOT_DIR: DOTENV_LOCAL_ROOT } = dotenv.parse(
      fs.readFileSync(dotenvLocalFile),
    );
    localIsRoot = DOTENV_LOCAL_ROOT;
  }

  if (fs.existsSync(dotenvFile)) {
    const { IS_PROJECT_ROOT_DIR: DOTENV_ROOT } = dotenv.parse(fs.readFileSync(dotenvFile));
    isRoot = DOTENV_ROOT;
  }

  return localIsRoot !== undefined ? localIsRoot !== 'false' : isRoot !== 'false';
};

/**
 * Return tsconfig compilerOptions.
 *
 * @param {string} directory
 * @returns {object}
 */
const getTsCompilerOptions = (directory) => {
  const tsconfigFile = path.resolve(directory, './tsconfig.json');
  let tsCompilerOptions = {};

  if (fs.existsSync(tsconfigFile)) {
    const { compilerOptions = { outDir: './dist', baseUrl: './src' } } = require(tsconfigFile);
    tsCompilerOptions = compilerOptions;
  }

  return tsCompilerOptions;
};

/**
 * Setup a webpack dotenv plugin config.
 *
 * @param {string} path
 * @returns {*}
 */
const setupWebpackDotenvFile = (path) => {
  const settings = {
    systemvars: true,
    silent: true,
  };

  if (path) {
    settings.path = path;
  }

  return new Dotenv(settings);
};

/**
 * Setup multiple webpack dotenv file parameters.
 *
 * @param {string} directory
 * @param {string} env
 * @param {boolean} isRoot
 * @returns {Array}
 */
const setupWebpackDotenvFilesForEnv = ({ directory, env, isRoot = true }) => {
  const dotenvWebpackSettings = [];

  if (env) {
    dotenvWebpackSettings.push(
      setupWebpackDotenvFile(path.resolve(directory, `.env.${env}.local`)),
    );
    dotenvWebpackSettings.push(setupWebpackDotenvFile(path.resolve(directory, `.env.${env}`)));
  }

  dotenvWebpackSettings.push(setupWebpackDotenvFile(path.resolve(directory, '.env.local')));
  dotenvWebpackSettings.push(setupWebpackDotenvFile(path.resolve(directory, '.env')));

  if (!isRoot) {
    if (env) {
      dotenvWebpackSettings.push(
        setupWebpackDotenvFile(path.resolve(directory, '..', `.env.${env}.local`)),
      );
      dotenvWebpackSettings.push(
        setupWebpackDotenvFile(path.resolve(directory, '..', `.env.${env}`)),
      );
    }

    dotenvWebpackSettings.push(setupWebpackDotenvFile(path.resolve(directory, '..', '.env.local')));
    dotenvWebpackSettings.push(setupWebpackDotenvFile(path.resolve(directory, '..', '.env')));
  }

  return dotenvWebpackSettings;
};

/**
 * Setup, and access, a dotenv file and the related set of parameters.
 *
 * @param {string} path
 * @returns {*}
 */
const setupDotenvFile = (path) => {
  const dotenvInitial = dotenv.config({ path });
  dotenvExpand(dotenvInitial);
};

/**
 * Setup and access local and specific dotenv file parameters.
 *
 * @param {string} env
 */
const setupDotenvFilesForEnv = ({ env }) => {
  const RELATIVE_DIRNAME = path.resolve(__dirname, '..');
  const IS_ROOT = getProjectIsRootDir(RELATIVE_DIRNAME);
  const { baseUrl: TS_BASE_URL, outDir: TS_OUT_DIR } = getTsCompilerOptions(RELATIVE_DIRNAME);

  if (!IS_ROOT) {
    if (env) {
      setupDotenvFile(path.resolve(RELATIVE_DIRNAME, '..', `.env.${env}.local`));
      setupDotenvFile(path.resolve(RELATIVE_DIRNAME, '..', `.env.${env}`));
    }

    setupDotenvFile(path.resolve(RELATIVE_DIRNAME, '..', '.env.local'));
    setupDotenvFile(path.resolve(RELATIVE_DIRNAME, '..', '.env'));
  }

  if (env) {
    setupDotenvFile(path.resolve(RELATIVE_DIRNAME, `.env.${env}.local`));
    setupDotenvFile(path.resolve(RELATIVE_DIRNAME, `.env.${env}`));
  }

  setupDotenvFile(path.resolve(RELATIVE_DIRNAME, '.env.local'));
  setupDotenvFile(path.resolve(RELATIVE_DIRNAME, '.env'));

  const DEPLOYMENT_MODE = process.env.DEPLOYMENT_MODE || 'integrated';
  const IMAGES_DIRNAME = process.env.IMAGES_DIRNAME || 'images';
  const PUBLIC_PATH = process.env.PUBLIC_PATH || '/';
  const SRC_DIR = path.resolve(RELATIVE_DIRNAME, process.env.SRC_DIR || TS_BASE_URL || 'src');
  const COMMON_DIR = path.resolve(RELATIVE_DIRNAME, process.env.COMMON_DIR || '../common');
  const DIST_DIR = path.resolve(RELATIVE_DIRNAME, process.env.DIST_DIR || TS_OUT_DIR || 'public');
  const HOST = process.env.HOST || DEPLOYMENT_MODE === 'integrated' ? '0.0.0.0' : 'localhost';
  const PORT = process.env.PORT || '9000';
  const PROXY_PROTOCOL = process.env.PROXY_PROTOCOL || 'http';
  const PROXY_HOST = process.env.PROXY_HOST || 'localhost';
  const PROXY_PORT = process.env.PROXY_PORT || process.env.PORT || 4000;
  const DEV_MODE = process.env.DEV_MODE || undefined;
  const OUTPUT_ONLY = process.env._OUTPUT_ONLY === 'true';

  process.env._RELATIVE_DIRNAME = RELATIVE_DIRNAME;
  process.env._IS_PROJECT_ROOT_DIR = IS_ROOT;
  process.env._IMAGES_DIRNAME = IMAGES_DIRNAME;
  process.env._PUBLIC_PATH = PUBLIC_PATH;
  process.env._SRC_DIR = SRC_DIR;
  process.env._COMMON_DIR = COMMON_DIR;
  process.env._DIST_DIR = DIST_DIR;
  process.env._HOST = HOST;
  process.env._PORT = PORT;
  process.env._PROXY_PROTOCOL = PROXY_PROTOCOL;
  process.env._PROXY_HOST = PROXY_HOST;
  process.env._PROXY_PORT = PROXY_PORT;
  process.env._OUTPUT_ONLY = OUTPUT_ONLY;
  process.env._DEV_MODE = DEV_MODE;
  process.env._DEPLOYMENT_MODE = DEPLOYMENT_MODE;
};

module.exports = { setupWebpackDotenvFilesForEnv, setupDotenvFilesForEnv };
