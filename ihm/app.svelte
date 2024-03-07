<!--permet d'appeler le composant svelte dans le fichier html avec une simple balise-->
<svelte:options customElement="app-todo" />

<script>
	let todos = [
		{ done: false, text: 'Ma 1ère tâche' },
		{ done: false, text: 'Ma 2ème tâche' },
		{ done: false, text: '' }
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
	h1{
		font-size: 70px;
	}

	p{
		font-size: large;
		margin-left: 20%;
		margin-right: auto;
	}
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
		font-size: large;
	}

	button{
		background-color: #216fedd3;
		width: 125px;
		border: none;
		color: white;
		padding: 15px 32px;
		text-align: center;
		text-decoration: none;
		display: inline-block;
		font-size: 16px;
	}

	button:hover{
		background-color: #0a47a9d3;
		cursor: pointer;
	}
	
</style>
