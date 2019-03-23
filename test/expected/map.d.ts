export type Person = {
  id?: string
  name?: string
}

export type PeopleById = {
  by_id?: Map<string, {
      id?: string
      name?: string
    }>
}