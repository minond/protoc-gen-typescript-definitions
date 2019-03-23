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

export type NestedPeople = {
  person1?: {
    id?: string
    name?: string
  }
  person2?: {
    id?: string
    name?: string
  }
  by_id?: {
    by_id?: Map<string, {
      id?: string
      name?: string
    }>
  }
}

export type NestedPeopleMap = {
  by_id?: Map<string, {
    person1?: {
      id?: string
      name?: string
    }
    person2?: {
      id?: string
      name?: string
    }
    by_id?: {
      by_id?: Map<string, {
        id?: string
        name?: string
      }>
    }
  }>
}