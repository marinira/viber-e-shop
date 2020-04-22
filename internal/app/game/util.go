package game

import "errors"

//функция вібирает из JSON объекта нужные значения в глубь дерева согласно названия нод
func Clash(myObj interface{}, a ...string) (interface{}, error) {
	//коныертируе входной аргумент в мап
	node1, ok := myObj.(map[string]interface{})
	if !ok {
		return nil, errors.New("Clash: Error convert interface to map ")
	}

	var ret interface{}
	//проходимся по названиям нод и по очереди каждую ноду конвертируем в мап
	//и дак до предпоследней ноды
	for i, v := range a {
		if i < len(a)-1 {
			node1, ok = node1[v].(map[string]interface{})
			if !ok {
				return nil, errors.New("Clash: Error convert interface to map ")
			}
		}
	}
	//вытаскиваем последнюю ноду по названию как интерфейс
	ret, ok = node1[a[len(a)-1]].(interface{})
	if !ok {
		return nil, errors.New("Clash: Error convert interface to map ")
	}
	return ret, nil
}
