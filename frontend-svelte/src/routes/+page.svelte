<script>
	let activeJobs = $state([{ id: 'API_WAIT', type: 'System', status: 'Standby' }]);
	let isUploading = $state(false);

	const fetchJobs = async () => {
		try {
			const res = await fetch('/api/v1/jobs');
			if (res.ok) {
				const data = await res.json();
				if (data && data.length > 0) {
					activeJobs = data;
				} else {
					activeJobs = [{ id: 'API_WAIT', type: 'System', status: 'Standby', progress: 0 }];
				}
			}
		} catch (err) {}
	};

	$effect(() => {
		fetchJobs();
		const interval = setInterval(fetchJobs, 3000);
		return () => clearInterval(interval);
	});

	const clearHistory = async () => { await fetch('/api/v1/jobs', { method: 'DELETE' }); fetchJobs(); };

	const uploadFile = async (file) => {
		if (!file) return;
		isUploading = true;
		const formData = new FormData();
		formData.append('document', file);
		try {
			await fetch('/api/v1/documents', {
				method: 'POST',
				body: formData
			});
		} catch (err) {}
		isUploading = false;
	};

	const handleDragOver = (e) => { e.preventDefault(); };
	const handleDrop = async (e) => {
		e.preventDefault();
		const file = e.dataTransfer.files[0];
		if (!file) return;
		const formData = new FormData();
		formData.append('document', file);
		await fetch('/api/v1/documents', { method: 'POST', body: formData });
		console.log("Upload triggered for:", file.name);
	};
	const handleFileInput = (e) => {
		if (e.target && e.target.files) {
			uploadFile(e.target.files[0]);
			e.target.value = '';
		}
	};
</script>

<svelte:window ondragover={(e) => e.preventDefault()} ondrop={(e) => e.preventDefault()} />

<div class="p-8 max-w-7xl mx-auto">
	<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
		<div class="lg:col-span-1">
			<label 
				ondragenter={handleDragOver} 
				ondragleave={handleDragOver} 
				ondragover={handleDragOver} 
				ondrop={handleDrop} 
				class="bg-slate-900/40 backdrop-blur-lg border border-white/10 shadow-xl rounded-2xl p-6 h-64 flex flex-col items-center justify-center border-dashed hover:bg-slate-800/40 transition-colors cursor-pointer"
			>
				<div class="text-4xl mb-4">📄</div>
				<div class="font-medium">Drag & Drop Files Here</div>
				<div class="text-sm text-slate-400 mt-2">PDF, DOCX, PNG up to 50MB</div>
				<input type="file" class="hidden" onchange={handleFileInput} />
			</label>
			
			<div class="bg-slate-900/40 backdrop-blur-lg border border-white/10 shadow-xl rounded-2xl p-6 mt-6">
				<h3 class="text-lg font-semibold mb-2">Token Wallet</h3>
				<div class="text-3xl font-bold text-blue-400">12,450</div>
				<div class="text-sm text-slate-400">Credits Remaining</div>
				<button class="w-full mt-4 bg-blue-600 hover:bg-blue-700 text-white py-2 rounded-lg transition-colors">
					Refill Wallet
				</button>
			</div>
		</div>
		
		<div class="lg:col-span-2">
			<div class="bg-slate-900/40 backdrop-blur-lg border border-white/10 shadow-xl rounded-2xl p-6 min-h-full">
				<div class="flex justify-between items-center mb-6">
					<h3 class="text-lg font-semibold">Active Jobs</h3>
					<button onclick={clearHistory} class="text-xs bg-red-500/10 text-red-400 border border-red-500/20 px-3 py-1.5 rounded-md hover:bg-red-500/20 transition-colors">Clear History</button>
				</div>
				<div class="space-y-4">
					{#each activeJobs as job}
						<div class="bg-slate-800/50 rounded-xl p-4 border border-white/5">
							<div class="flex justify-between items-center mb-2">
								<div class="font-medium">
									{job.type} <span class="text-xs text-slate-400 ml-2">{job.id}</span>
								</div>
								<div class="text-sm px-2 py-1 rounded-md capitalize {job.status === 'PROCESSING' ? 'bg-blue-500/20 text-blue-300' : 'bg-green-500/20 text-green-300'}">
									{job.status}
								</div>
							</div>
							<div class="h-2 w-full bg-slate-700 rounded-full overflow-hidden">
								<div class="h-full bg-blue-500 transition-all duration-500" style="width:100%"></div>
							</div>
							{#if job.res}
								<div class="mt-3 p-4 bg-slate-950/80 rounded-lg overflow-x-auto text-[11px] text-emerald-400 font-mono border border-emerald-500/20 shadow-inner">
									<pre>{JSON.stringify(job.res, null, 2)}</pre>
								</div>
							{/if}
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>
</div>
