<svelte:options customElement="signin-todo" />

<script>

    let email = '';
    let password = '';
    let confirmpassword = '';
  
    const handleSubmit = () => {
        let newuser = {
            email,
            password,
            confirmpassword
        }
        console.log("Email:", email);
        console.log("Password:", password);
        console.log("Password2:", confirmpassword);
        try {
            sendUser(newuser);
        } catch (error) {
			console.error(`Erreur lors de la connection au serveur : ${error.message}`);
		}
    };

    function customQueryEscape(str){
		const uriStr = encodeURIComponent(str);
		const finalQueryStr = uriStr
		.replace(/!/g, '%21')
		.replace(/'/g, '%27')
		.replace(/\(/g, '%28')
        .replace(/\)/g, '%29')
		.replace(/\*/g, '%2A')
		.replace(/%20/g,'+');
    	return finalQueryStr;
	}

    async function sendUser(user) {
		user.email = customQueryEscape(user.email);
        user.password = customQueryEscape(user.password);
        user.confirmpassword = customQueryEscape(user.confirmpassword);
        const url =`/signin?email=${user.email}&password=${user.password}&confirmpassword=${user.confirmpassword}`;
		try {
			const reponse = await fetch(url,{method: "POST"});
			if (!reponse.ok) {
				const errorData = await reponse.text();
				alert(errorData);
				throw new Error(`Erreur lors de la requête : ${reponse.status} ${reponse.statusText}`);
			}
			await reponse.json();
		} catch (error) {
			console.error(`Une erreur s'est produite : ${error.message}`);
		}
	}
    

</script>

<style>

    h1{
		font-size: 70px;
	}

	label{
		font-size: large;
	}
    .centered {
		width: 30em;
		margin: auto;
		display: inline-flexbox;
	}
</style>

<div class="centered">

    <h1>My TodoList</h1>
    
    <form on:submit|preventDefault={handleSubmit}>
        <label>
        Email:
            <input type="email" bind:value={email} required>
        </label>
        <br>
        <label>
        Mot de passe:
            <input type="password" bind:value={password} required>
        </label>
        <br>
        <label>
        Confirmer le mot de passe:
            <input type="password" bind:value={confirmpassword} required>
        </label>
        <br>
        <button type="submit">S'inscrire</button>
    </form>
</div>