import * as esbuild from 'esbuild';
import sveltePlugin from 'esbuild-svelte';

let ctx = await esbuild.context({
	entryPoints: ['app.svelte'],
	bundle: true,
	format: 'esm',
	outdir: './build',
	plugins: [
		sveltePlugin({
			compilerOptions: { customElement: true}
		})
	],
	banner: {
		js: "new EventSource('http://127.0.0.1:8888/esbuild').addEventListener('change', () => location.reload())"
	},
	logLevel: 'info'
});
await ctx.watch();
await ctx.serve({
	servedir: './',
	port: 8888,
	host: '127.0.0.1'
});