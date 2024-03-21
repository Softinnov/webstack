<!--permet d'appeler le composant svelte dans le fichier html avec une simple balise-->
<svelte:options customElement="app-todo" />

<script>
	let todos = [];
	let nouvelleTache='';

	function add() {
		let todo = {
			action: "add",
			text: nouvelleTache
		}
		try{
			sendTodo(todo);
		} catch (error) {
			console.error(`Erreur lors de la connection au serveur : ${error.message}`);
		}
		nouvelleTache='';
	}

	function xclear(item) {
		try {
			item.action="delete";
			sendTodo(item);
		} catch (error) {
			console.error(`Erreur lors de la connection au serveur : ${error.message}`);
		}
	}

	function modify(item) {
		try {
			item.action="modify";
			sendTodo(item);
		} catch (error) {
			console.error(`Erreur lors de la connection au serveur : ${error.message}`);
		}
	}

	async function getTodoList() {
		const url = `/todos`	
		try {
			const reponse = await fetch(url,{method: "GET"});
			if (!reponse.ok) {
				const errorData = await reponse.text();
				alert(errorData);
				throw new Error(`Erreur lors de la requête : ${reponse.status} ${reponse.statusText}`);
			}
			const result = await reponse.json();
			if(result != null){
				todos = result;
			}
			if(result == null){
				todos = [];
			}
			console.log(todos);
		} catch (error) {
			console.error(`Une erreur s'est produite : ${error.message}`);
		}
	}

	async function sendTodo(todo) {
		const url = `/service?id=${todo.id}&action=${todo.action}&text=${todo.text}`
		
		try {
			const reponse = await fetch(url,{method: "POST"});
			if (!reponse.ok) {
				const errorData = await reponse.text();
				alert(errorData);
				throw new Error(`Erreur lors de la requête : ${reponse.status} ${reponse.statusText}`);
			}
			const result = await reponse.json();
			getTodoList();
		} catch (error) {
			console.error(`Une erreur s'est produite : ${error.message}`);
		}
	}

	function handleKeydown(e) {
		if(e.key=="Enter") {
			add();
		}
	}

	$: remaining = todos.length;	
	
	getTodoList();

</script>

<div class="centered">

	<h1>My TodoList</h1>

	<input 
		class="newTask"
		type="text" 
		placeholder="Quoi d'autre?"
		bind:value={nouvelleTache}
		on:keydown={handleKeydown}
	/>
	<button class="ajout" disabled={nouvelleTache==""} on:click={add}>
		✔️
	</button>
	<ul id="todo-list" class="todos">	
		{#each todos as item}
			<li class:done={item.done}>
				<button type="button" class="button" on:click={modify(item)}>
					✏️
				</button>
				<button type="button" class="button" on:click={xclear(item)}>
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
	.button {
		height: 30px;
		width: 30px;
		border: none;
		border-radius: 5%;
		background-color: white;
		margin-right: 2%;
	}
	.button:hover{
		background-color: rgba(146, 146, 146, 0.381);
		cursor: pointer;
	}
	.centered {
		max-width: 20em;
		margin: 0 auto;
	}
	.newTask {
		margin-bottom: 15%;
		margin-right: 1%;
	}
	.todos {
		margin-top: 10%;
	}
	.ajout {
		height: 30px;
		border: none;
		border-radius: 5%;
		background-color: white;
		position: relative;
		right: -5%;
		top: -10%
	}
	.ajout:hover{
		background-color: rgba(146, 146, 146, 0.381);
		cursor: pointer;
	}
	.ajout:disabled{
		background-color: white;
		color: rgba(128, 128, 128, 0.836);
		cursor: default;
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
