<!--permet d'appeler le composant svelte dans le fichier html avec une simple balise-->
<svelte:options customElement="app-todo" />

<script>

	let todos = [
		{ done: false, text: 'Ma 1ère tâche' }
	];
	let nouvelleTache='';

	function add() {
		let todo = {
			done: false,
			text: nouvelleTache
		}
		
		try{
			reqFetch(todo);
			alert("ToDo bien ajouté !");
		} catch (error) {
			console.error(`Erreur lors de la connection au serveur : ${error.message}`)
		}
		todos = todos.concat(todo);
		nouvelleTache='';		
	}

	function clear() {
		

		todos = todos.filter((t) => !t.done);
	}

	async function reqFetch(todo) {
		const url = `/service?check=${todo.done}&text=${todo.text}`
		
		try {
			const reponse = await fetch(url,{method: "GET"});
			console.log(reponse);
			if (!reponse.ok) {
				throw new Error(`Erreur lors de la requête : ${reponse.status} ${reponse.statusText}`);
			}
			const result = await reponse.json();
			console.log(result);
		} catch (error) {
			console.error(`Une erreur s'est produite : ${error.message}`);
		}
	}

	$: remaining = todos.filter((t) => !t.done).length;
	$: deletable = todos.filter((t)=>t.done).length>0;
</script>

<div class="centered">
	<h1>Ma ToDoList</h1>

	{JSON.stringify(todos)}
	<input 
		type="text" 
		placeholder="Quoi d'autre?"
		bind:value={nouvelleTache}
	/>
	<button class="ajout" disabled={nouvelleTache==""} on:click={add}>
		Ajouter
	</button>
	<form method="GET" class="todos">	
		{#each todos as todo}
			<li class:done={todo.done}>
				<input
					id="check"
					type="checkbox"
					bind:checked={todo.done}
				/>

				<input
					id="todo"
					type="text"
					bind:value={todo.text}
				/>
			</li>
		{/each}
		</form>

	<p>{remaining} tâches restantes !</p>

	<button disabled={!deletable} on:click={clear}>
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

	.ajout {
		position: relative;
		right: -85%;
		top: -10%
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
	
</style>
