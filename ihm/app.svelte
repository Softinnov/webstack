<!--permet d'appeler le composant svelte dans le fichier html avec une simple balise-->
<svelte:options customElement="app-todo" />

<script>
	let todos = [
		{ done: false, text: 'Ma 1ère tâche' },
		{ done: false, text: 'Ma 2ème tâche' }
	];

	function add() {
		todos = todos.concat({
			done: false,
			text: ''
		});
	}

	function clear() {
		todos = todos.filter((t) => !t.done);
	}

	$: remaining = todos.filter((t) => !t.done).length;
</script>

<div class="centered">
	<h1>Ma ToDoList</h1>

	<ul class="todos">
		{#each todos as todo}
			<li class:done={todo.done}>
				<input
					type="checkbox"
					bind:checked={todo.done}
				/>

				<input
					type="text"
					placeholder="Quoi d'autre?"
					bind:value={todo.text}
				/>
			</li>
		{/each}
	</ul>

	<p>{remaining} tâches restantes !</p>

	<button on:click={add}>
		Ajouter
	</button>

	<button on:click={clear}>
		Supprimer
	</button>
</div>

<style>
	.centered {
		max-width: 20em;
		margin: 0 auto;
	}

	.done {
		opacity: 0.4;
	}

	li {
		display: flex;
	}

	input[type="text"] {
		flex: 1;
		padding: 0.5em;
		margin: -0.2em 0;
		border: none;
	}
</style>
