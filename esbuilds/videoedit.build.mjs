import * as esbuild from 'esbuild';

await esbuild.build({
  entryPoints: ['ts/videoedit.page.ts'],
  bundle: true,
  minify: true,
  outfile: 'templates/js/videoedit.page.js',
});