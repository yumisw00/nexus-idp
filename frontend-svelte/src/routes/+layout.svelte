<script>import '../app.css';</script>
<div class="min-h-screen text-slate-100 font-sans relative">
<div class="fixed inset-0 bg-slate-900 -z-10 overflow-hidden">
<div class="absolute top-10 left-10 w-96 h-96 bg-blue-600/30 blur-[120px] mix-blend-multiply opacity-70 rounded-full"></div>
<div class="absolute bottom-10 right-10 w-96 h-96 bg-purple-600/20 blur-[120px] mix-blend-multiply opacity-70 rounded-full"></div>
</div>
<nav class="sticky top-0 backdrop-blur-md bg-slate-900/40 border-b border-white/10 p-4"><div class="container mx-auto font-bold text-xl">Nexus IDP</div></nav>
<main class="container mx-auto p-4"><slot/></main>
</div>
