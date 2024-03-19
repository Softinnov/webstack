import * as esbuild from 'esbuild';
import sveltePlugin from 'esbuild-svelte';

const isDev = process.argv.includes('--dev')

let commonOptions = {
	entryPoints: ['ihm/app.svelte'],
	bundle: true,
	format: 'esm',
	plugins: [
		sveltePlugin({
			compilerOptions: { customElement: true}
		})
	],
}
let devOptions = {
	...commonOptions,
	outdir: './ihm',
	banner: {
        //crée un eventlistener qui détecte les modifications du code et actualise la page pour afficher les modifs en direct
		js: "new EventSource('http://127.0.0.1:8888/esbuild').addEventListener('change', () => location.reload())"
	},
	logLevel: 'info'
};
let prodOptions = {
	...commonOptions,
	outdir: './build',
	minify: true,
	logLevel: 'error'
}

if (!isDev) {
	await esbuild.build(prodOptions);
	process.exit(0);
}

let ctx = await esbuild.context(devOptions);
await ctx.watch();
await ctx.serve({
	servedir: './',
	port: 8888,
	host: '127.0.0.1'
});
