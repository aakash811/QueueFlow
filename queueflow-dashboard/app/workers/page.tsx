export default function WorkersPage() {
  return (
    <div>
      <h1 className="text-4xl font-bold mb-8">Workers</h1>

      <div className="grid grid-cols-4 gap-6">
        {[1, 2, 3, 4].map((worker) => (
          <div
            key={worker}
            className="bg-zinc-900 p-6 rounded-2xl border border-zinc-800"
          >
            <h2 className="text-xl font-semibold">Worker {worker}</h2>

            <p className="text-green-400 mt-3">Healthy</p>
          </div>
        ))}
      </div>
    </div>
  );
}
