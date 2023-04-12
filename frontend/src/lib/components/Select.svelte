<script lang="ts">
    export let options: Map<string, string>;
    export let value: string = "default";
    export let name: string;
    export let placeholder: string;
</script>

<label for={name} class="block mb-2 text-sm font-medium text-gray-900 dark:text-gray-400">{placeholder}</label>
<select bind:value id={name} class="mb-5 bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
  <option selected={value == "default"} value={"default"}>{placeholder}</option>
  {#each [...options] as option}
    <option value={option[1]} selected={option[1] == value}>{option[0]}</option>
  {/each}
</select>