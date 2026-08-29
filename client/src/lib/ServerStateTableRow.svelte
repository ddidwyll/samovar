<script>
  let { name, value, unit } = $props()
</script>

<tr>
  <th>{name}:</th>
  <td>
    {value}
    {unit ?? ""}
  </td>
</tr>

<style>
  th, td {
    font-weight: normal;
    padding: 10px 20px;
    border: solid 2px #aaa;
    min-width: 40px;
  }
  th {
    text-align: start;
    text-transform: uppercase;
  }
  td {
    text-align: center;
  }
</style>
