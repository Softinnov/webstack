<!--permet d'appeler le composant svelte dans le fichier html avec une simple balise-->
<svelte:options customElement="app-todo" />

<script>
	let todos = [];
	let nouvelleTache='';

	function add() {
		let todo = {
			done: false,
			text: nouvelleTache
		}
		try{
			sendTodo(todo);
		} catch (error) {
			console.error(`Erreur lors de la connection au serveur : ${error.message}`);
		}
		getTodoList();
		nouvelleTache='';
	}

	function xclear(item) {
		try {
			item.done=true;
			sendTodo(item);
		} catch (error) {
			console.error(`Erreur lors de la connection au serveur : ${error.message}`);
		}
		getTodoList();
		// todos = todos.filter((t) => t.done==false);
	}

	async function getTodoList() {
		const url = `/todos`	
		try {
			const reponse = await fetch(url,{method: "GET"});
			if (!reponse.ok) {
				// const errorData = await reponse.text();
				// alert(errorData.message);
				throw new Error(`Erreur lors de la requête : ${reponse.status} ${reponse.statusText}`);
			}
			const result = await reponse.json();
			console.log(result);
			todos = result;
			
		} catch (error) {
			console.error(`Une erreur s'est produite : ${error.message}`);
		}
	}

	async function sendTodo(todo) {
		const url = `/service?check=${todo.done}&text=${todo.text}`
		
		try {
			const reponse = await fetch(url,{method: "POST"});
			if (!reponse.ok) {
				const errorData = await reponse.text();
				alert(errorData);
				throw new Error(`Erreur lors de la requête : ${reponse.status} ${reponse.statusText}`);
			}
			const result = await reponse.json();
			// console.log(result);
		} catch (error) {
			console.error(`Une erreur s'est produite : ${error.message}`);
		}
	}

	$: remaining = todos.filter((t) => !t.done).length;

</script>

<div class="centered">

	<h1>Ma ToDoList</h1>

	<input 
		class="newTask"
		type="text" 
		placeholder="Quoi d'autre?"
		bind:value={nouvelleTache}
	/>
	<button class="ajout" disabled={nouvelleTache==""} on:click={add}>
		Ajouter
	</button>
	<ul id="todo-list" class="todos">	
		{#each todos as item}
			<li class:done={item.done}>
				<button type="button" class="delete" on:click={xclear(item)}>
					🗑️
				</button>

				<input
					id="todo"
					type="text"
					bind:value={item.text}
				/>
			</li>
		{/each}
		</ul>

	<p>{remaining} tâches restantes !</p>
	
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

	.newTask {
		margin-bottom: 15%;
	}

	.todos {
		margin-top: 10%;
	}

	.ajout {
		position: relative;
		right: -15%;
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
