package main

import "fmt"

type Person struct {
	name string
	addrInfo map[string]string
}

var PersonList = make([]Person,0)

func ScanNum(){
	var input int
	fmt.Println("添加联系人信息，请按1")
	fmt.Println("删除联系人信息，请按2")
	fmt.Println("查询联系人信息，请按3")
	fmt.Println("编辑联系人信息，请按4")
	fmt.Scan(&input)
	Switctype(input)

}

func ShowPerons(){
	if len(PersonList) > 0{
		for _,val := range PersonList{
			fmt.Println("姓名：", val.name)
			for k, v := range val.addrInfo {
				fmt.Println("电话类型：", k)
				fmt.Println("电话号码：", v)
			}
		}
	}else{
		fmt.Println("暂时没有联系人信息")
	}
}
func AddPerson(){
	var name string
	var phoneType string
	var phoneNum string

	var addrPhone = make(map[string]string)
	var keyInput string
	fmt.Println("请输入姓名")
	fmt.Scan(&name)
	for{
		fmt.Println("请输入电话类型")
		fmt.Scan(&phoneType)

		fmt.Println("请输入电话号码")
		fmt.Scan(&phoneNum)

		addrPhone[phoneType] = phoneNum

		fmt.Println("如果结束电话的录入，请按Q")
		fmt.Scan(&keyInput)
		if keyInput == "Q"{
			break
		}
	}

	PersonList = append(PersonList,Person{name:name,addrInfo:addrPhone})
	ShowPerons()
}

func DelPerson(){
	var name string
	var index = -1
	fmt.Println("请输入要删除的联系人姓名：")
	fmt.Scan(&name)

	for k,v := range PersonList {
		if(v.name == name){
			index = k
			break;
		}
	}

	if(index != -1){
		PersonList = append(PersonList[:index],PersonList[index+1:]...)
	}
	ShowPerons()
}

func FindPerson() *Person{
	var name string
	var index = -1
	fmt.Println("请输入要查询的联系人姓名：")
	fmt.Scan(&name)

	for k,v := range PersonList {
		if(v.name == name){
			index = k
			fmt.Println("联系人姓名：", v.name)
			for t, phone := range v.addrInfo{
				fmt.Printf("%s:%s\n", t, phone)
			}
			break;
		}
	}

	if(index == -1){
		fmt.Println("没有找到联系人信息")
		return nil
	}else{
		return &PersonList[index]
	}

}

func EditPerson(){
	var p *Person = nil
	var num int = 0
	var name string
	var phonetType int = -1
	var phoneNum string
	p = FindPerson()
	menu := make([]string,0)

	if(p != nil) {
		for {
			fmt.Println("编辑用户名称请按:5,编辑电话请按:6,退出请按:7")
			fmt.Scan(&num)

			switch  num {
			case 5:
				fmt.Println("请输入新的姓名：")
				fmt.Scan(&name)
				p.name = name
				ShowPerons()

			case 6:
				var j int
				for key, value := range p.addrInfo {
					fmt.Println("编辑(", key, ")", value, "请按：", j)
					menu = append(menu, key)
					j++
				}

				fmt.Println("请输入编辑号码的类型：")
				fmt.Scan(&phonetType)
				fmt.Println("请输入新的电话号码：")
				fmt.Scan(&phoneNum)
				p.addrInfo[menu[phonetType]] = phoneNum

			case 7:
				break
			default:
				continue
			}
			if num == 7{
				break
			}
		}
	}else{
		fmt.Println("没有找到要编辑的联系人信息")
	}

}

func Switctype(typeNum int){
	switch typeNum {
	case 1:
		AddPerson();
	case 2:
		DelPerson();
	case 3:
		FindPerson();
	case 4:
		EditPerson();
	default:
		fmt.Println("Invalid input type")
	}
}

func main(){
	for {
		ScanNum()
	}

}
