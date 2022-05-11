export default function tableMenuItems(dictionary) {
  return [
    {
      name: "deleteTable",
      tooltip: dictionary.deleteTable,
      icon: "editor-delete_table",
      active: () => false,
    },
  ];
}
