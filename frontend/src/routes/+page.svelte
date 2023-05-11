<script>
	import CopyClipBoard from './CopyClipBoard.svelte';

	let progress = 0;
	let files, input;
	let xhr, xhrRes;
	let link;

	function uploadAnother() {
		// reset to zero / null
		progress = 0;
		xhr, files, xhrRes = null;
		input.value = '';
	}

	function cancelUpload() {
		if (xhr) {
			xhr.abort(); // abort xhr upload
		}
		uploadAnother();
	}

	const copy = () => {
		const app = new CopyClipBoard({
			target: document.getElementById('clipboard'),
			props: { value: link }
		});
		app.$destroy();
	};

    function uploadFile(files) {
		if (files.length > 0) {
			let file = files[0];
            if (Math.floor(file.size / 1000000) > 100) {
                alert('File size is too big. Max 100MB');
                uploadAnother()
                return;
            }

			var formdata = new FormData();
			formdata.append('file', file);

			xhr = new XMLHttpRequest();
			xhr.upload.onprogress = function (e) {
				progress = Math.ceil((e.loaded / e.total) * 100);
			};
			xhr.onreadystatechange = function () {
				if (xhr.readyState == XMLHttpRequest.DONE) {
					xhrRes = JSON.parse(xhr.responseText);
					link = xhrRes.link;
				}
			};
			xhr.open('POST', import.meta.env.VITE_BASE_URL+"/", true);
			xhr.send(formdata);
		}
	}

	$: if (files) {
		uploadFile(files);
	}
</script>

<svelte:head>
	<title>Temporar File Sharing</title>
</svelte:head>

<header id="header">
	<div class="content">
		<h1><a href="/">Temporar File Sharing</a></h1>
		<p>Ephemeral file sharing <br /> Convenient, anonymous and secure</p>
		<ul class="actions">
			<li class="upload-container">
                {#if progress > 0 && progress < 100}
				<div style="min-width: 330px;">
					<div class="progress-bar">
						<span class="progress-bar-fill" style="width: {progress}%;" />
					</div>
					{#if files}
						{files[0].name}
						{progress}% of {Math.floor(files[0].size / 1000000) + 'MB'}

						{#if xhr}
							<a on:click={cancelUpload} href="#/">Cancel</a>
						{/if}
					{/if}
				</div>
                {/if}

				{#if xhrRes}
					<ul style="text-align: center;">
						<li class="alert alert-success">
							<div class="alert alert-success">
								<strong>Upload success.</strong> Here is the link to your file:
							</div>
							<span class="lead">
								<a target="_blank" href={xhrRes.link}>{xhrRes.link}</a>
								<button on:click={copy} style="vertical-align: top;">Copy Link</button>
								<button on:click={uploadAnother}>Upload Another</button>
								<br /><br />
							</span>
						</li>
					</ul>
				{:else}
					<ul class={progress != 100 ? 'hide' : ''}>
						<li class="alert alert-success">
							<div class="alert alert-success">Processing file...</div>
						</li>
					</ul>
				{/if}

                {#if progress == 0}
				<div class="button special">
					<div class="span12">
						<div
							class="btn btn-success"
							style="position: relative; overflow: hidden; direction: ltr;"
						>
							Upload
							<input bind:files bind:this={input} type="file" name="file" />
						</div>
					</div>
				</div>
                {/if}
			</li>
			<li><a href="#one" class="button">Learn More</a></li>
		</ul>
		<small>
            Fully inspired by <a href="https://file.io" target="_blank">file.io</a>, web content are
			copied from old file.io site <br> Powered by <a href="https://svelte.dev/" target="_blank">Svelte.dev</a>
        </small>
	</div>
</header>

<div id="clipboard" />

<section id="one" class="wrapper style2 special">
	<h2>"It's like snapchat, but for files!"</h2>
	- Ben A, Philadelphia
	<br />
	<br />
	<h4>
		Simply upload a file, share the link, and after it is downloaded,<br />
		the file is completely deleted. For added security, file are automatically <br />
		deleted after 2 hours, even if it was never downloaded.<br /><br />
		All files are (currently) not-encrypted when stored on our servers.<br /><br />
	</h4>
</section>

<!-- Two -->
<section id="api" class="wrapper">
	<div class="inner alt">
		<section class="spotlight">
			<div class="content">
				<h3>Easy to use API</h3>
				<p>
					Try it out:
					<br />
					<br />
					<code>
						$ curl -F "file=@test.txt" https://domain.com<br />
					</code>
				</p>
			</div>
		</section>
	</div>
</section>

<!-- Three -->
<section id="faq" class="wrapper">
	<div class="inner alt">
		<section class="spotlight">
			<div class="content">
				<h3>FAQ</h3>

				<h4>Are there log files or any backups of the file after it is deleted?</h4>
				<p>No, it is anonymous and we erase everything and There are no backups.</p>
				<h4>Is it free?</h4>
				<p>Yes!</p>
				<h4>What kinds of files can I share?</h4>
				<p>
					No illegal or copyrighted content is allowed. By using this service you agree to the <a
						href="/tos.html">Terms of Service</a
					>
				</p>
				<h4>Is there a size limit?</h4>
				<p>Yes, there is a 100MB per file limit.</p>
				<h4>Who are you and how can I trust you?</h4>
				<p>?</p>
			</div>
		</section>
	</div>
</section>

<footer id="footer">
	<p class="copyright">
		&copy; Copyright 2022. Design credits: <a href="http://html5up.net">HTML5 UP</a>
	</p>
</footer>

<style>
	.progress-bar {
		min-width: 100%;
		background-color: #e0e0e0;
		padding: 3px;
		border-radius: 3px;
		box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.2);
	}

	.progress-bar-fill {
		display: block;
		height: 22px;
		background-color: #659cef;
		border-radius: 3px;

		transition: width 500ms ease-in-out;
	}

	.hide {
		display: none;
	}

	input {
		position: absolute;
		right: 0px;
		top: 0px;
		font-family: Arial;
		font-size: 118px;
		margin: 0px;
		padding: 0px;
		cursor: pointer;
		opacity: 0;
	}
</style>
