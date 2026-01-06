import js from '@eslint/js';
import svelte from 'eslint-plugin-svelte';
import prettier from 'eslint-config-prettier';

/** @type {import('eslint').Linter.FlatConfig[]} */
export default [
	js.configs.recommended,
	...svelte.configs['flat/recommended'],
	prettier,
	...svelte.configs['flat/prettier'],
	{
		languageOptions: {
			globals: {
				// Browser globals
				window: 'readonly',
				document: 'readonly',
				navigator: 'readonly',
				console: 'readonly',
				localStorage: 'readonly',
				sessionStorage: 'readonly',
				fetch: 'readonly',
				URL: 'readonly',
				// Node globals
				process: 'readonly',
				__dirname: 'readonly',
				__filename: 'readonly'
			}
		}
	},
	{
		ignores: ['build/', '.svelte-kit/', 'dist/', 'node_modules/']
	}
];
