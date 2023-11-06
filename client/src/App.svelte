<script lang="ts">
  import "./app.css";
  import axios from "axios";

  let url: string;
  let isLoading: boolean = false;
  let shortenedUrl: string;

  const onSubmit = async () => {
    isLoading = true;

    const payload = {
      url: url,
    };

    // make a post request to the api
    const response = await axios.post("https://glitch-sazed.onrender.com/api/v1", payload);
    const body = response.data;

    shortenedUrl = body.short;

    isLoading = false;
  };
</script>

<main
  class="container mx-auto h-screen flex flex-col gap-4 justify-center items-center"
>
  <form
    class="w-full max-w-sm flex flex-col md:flex-row items-center gap-3"
    on:submit|preventDefault={onSubmit}
  >
    <label for="url" class="sr-only">Url</label>
    <input
      class="flex h-10 w-full rounded-md border border-gray-300 bg-transparent px-3 py-2 text-sm ring-offset-white file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-gray-600 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gray-300 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
      type="url"
      name="url"
      id="url"
      required
      placeholder="https://github.com/SazedWorldbringer"
      bind:value={url}
    />
    <button
      type="submit"
      class="w-full md:w-min inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 bg-gray-800 text-zinc-100 shadow hover:bg-gray-800/80 h-9 px-4 py-2"
      disabled={isLoading}
    >
      {#if isLoading}
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 20 20"
          fill="currentColor"
          class="mr-2 w-5 h-5 animate-spin"
        >
          <path
            fill-rule="evenodd"
            d="M15.312 11.424a5.5 5.5 0 01-9.201 2.466l-.312-.311h2.433a.75.75 0 000-1.5H3.989a.75.75 0 00-.75.75v4.242a.75.75 0 001.5 0v-2.43l.31.31a7 7 0 0011.712-3.138.75.75 0 00-1.449-.39zm1.23-3.723a.75.75 0 00.219-.53V2.929a.75.75 0 00-1.5 0V5.36l-.31-.31A7 7 0 003.239 8.188a.75.75 0 101.448.389A5.5 5.5 0 0113.89 6.11l.311.31h-2.432a.75.75 0 000 1.5h4.243a.75.75 0 00.53-.219z"
            clip-rule="evenodd"
          />
        </svg>
      {/if}
      Shorten
    </button>
  </form>

  <p>
    Here's your shortened URL: {#if shortenedUrl}
      <a target="_blank" href={`https://${shortenedUrl}`}>https://{shortenedUrl}</a>
    {:else}
      _
    {/if}
  </p>
</main>
