package entity

type Monster struct {
    Stats     Character
    Type      MonsterType
    Hostility HostilityType
    IsChasing bool
    Dir       Direction
}
