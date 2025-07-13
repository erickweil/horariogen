export const aulas = {
	"turmas": [
		{
			"nome": "Turma 2025",
			"horarios": {
				"dom":[],
				"seg":[1,2,3,4],
				"ter":[1,2,3,4],
				"qua":[1,2,3,4],
				"qui":[1,2,3,4],
				"sex":[1,2,3,4],
				"sab":[]
			}
		}, {
			"nome": "Turma 2024",
			"horarios": {
				"dom":[],
				"seg":[1,2,3,4],
				"ter":[1,2,3,4],
				"qua":[1,2,3,4],
				"qui":[1,2,3,4],
				"sex":[1,2,3,4],
				"sab":[]
			}
		}, {
			"nome": "Turma 2023",
			"horarios": {
				"dom":[],
				"seg":[1,2,3,4],
				"ter":[1,2,3,4],
				"qua":[1,2,3,4],
				"qui":[1,2,3,4],
				"sex":[1,2,3,4],
				"sab":[]
			}
		}
	],
	"disciplinas": [
		{
			"nome": "Fundamentos de Tecnologia da Informação",
			"turma": "Turma 2025",
			"aulas": 2,
			"agrupar": 2
		},
		{
			"nome": "Algoritmos e Lógica de Programação",
			"turma": "Turma 2025",
			"aulas": 4,
			"agrupar": 2
		},
		{
			"nome": "Matemática Computacional",
			"turma":"Turma 2025",
			"aulas": 4,
			"agrupar": 2
		},
		{
			"nome": "Organização de computadores",
			"turma":"Turma 2025",
			"aulas": 2,
			"agrupar": 2
		},
		{
			"nome": "Relações Étnico-Raciais",
			"turma":"Turma 2025",
			"aulas": 2,
			"agrupar": 2
		},
		{
			"nome": "Fundamentos em Negócios",
			"turma":"Turma 2025",
			"aulas": 2,
			"agrupar": 2
		},
		{
			"nome": "Banco de dados I",
			"turma":"Turma 2025",
			"aulas": 4
			
		},


		{
			"nome": "Estrutura de dados I",
			"turma":"Turma 2024",
			"aulas": 2,
			"agrupar": 2
		},
		{ 
			"nome":"Orientação a objetos",
			"turma": "Turma 2024",
			"aulas": 4,
			"agrupar": 2
		},
		{ 
			"nome":"Devops e Cloud I",
			"turma": "Turma 2024",
			"aulas": 2,
			"agrupar": 2
		},
		{ 
			"nome":"Teste de Software I",
			"turma": "Turma 2024",
			"aulas": 2,
			"agrupar": 2
		},
		{ 
			"nome":"Programação Web: front-end",
			"turma": "Turma 2024",
			"aulas": 4,
			"agrupar": 2
		},
		{ 
			"nome":"Programação Web: back-end",
			"turma": "Turma 2024",
			"aulas": 2,
			"agrupar": 2
		},
		{ 
			"nome":"Fábrica de software II",
			"turma": "Turma 2024",
			"aulas": 4,
			"agrupar": 4
		},

		{
			"nome": "Ciência de dados",
			"turma": "Turma 2023",
			"aulas": 2,
			"agrupar": 2
		},
		{
			"nome": "Tópicos Especiais I",
			"turma": "Turma 2023",
			"aulas": 2,
			"agrupar": 2
		},
		{
			"nome": "Metodologia Científica para Computação",
			"turma": "Turma 2023",
			"aulas": 2,
			"agrupar": 2
		},
		{
			"nome": "Internet das Coisas",
			"turma": "Turma 2023",
			"aulas": 4,
			"agrupar": 2
		},
		{
			"nome": "Dispositivos Móveis I",
			"turma": "Turma 2023",
			"aulas": 4,
			"agrupar": 2
		},
		{
			"nome": "Dispositivos Móveis II",
			"turma": "Turma 2023",
			"aulas": 2,
			"agrupar": 2
		},
		{
			"nome": "Fábrica de Software IV",
			"turma": "Turma 2023",
			"aulas": 4,
			"agrupar": 4
		}
	],
	"disciplinas_unidas": [
		{
			"grupo": "Fábrica de Software",
			"disciplinas": [
				"Fábrica de software II",
				"Fábrica de Software IV"
			]
		}
	],
	"professores": [
		{ 
			"nome":"Roberto",
			"disciplinas": [
				"Fundamentos de Tecnologia da Informação",
				"Organização de computadores",
				"Fábrica de software II",
				"Tópicos Especiais I",
				"Internet das Coisas"
			],
			"horarios": {
				"dom":[],
				"seg":[1,2,3,4],
				"ter":[1,2,3,4],
				"qua":[1,2,3,4],
				"qui":[],
				"sex":[1,2,3,4],
				"sab":[]
			}
		},
		{ 
			"nome":"Erick",
			"disciplinas": [
                "Algoritmos e Lógica de Programação",
				"Estrutura de dados I",
				"Devops e Cloud I",
				"Fábrica de software II"
			],
			"horarios": {
				"dom":[],
				"seg":[1,2,3,4],
				"ter":[1,2,3,4],
				"qua":[],
				"qui":[3,4],
				"sex":[1,2,3,4],
				"sab":[]
			}
		},
		{ 
			"nome":"Lucas",
			"disciplinas": [
                "Ciência de dados"
			],
			"horarios": {
				"dom":[],
				"seg":[1,2,3,4],
				"ter":[1,2,3,4],
				"qua":[1,2,3,4],
				"qui":[1,2,3,4],
				"sex":[],
				"sab":[]
			}
		},
		{ 
			"nome":"Wesley",
			"disciplinas": [
				"Banco de dados I",
				"Teste de Software I",
				"Fábrica de software II"
			],
			"horarios": {
				"dom":[],
				"seg":[],
				"ter":[1,2,3,4],
				"qua":[],
				"qui":[1,2,3,4],
				"sex":[1,2,3,4],
				"sab":[]
			}
		},
		{ 
			"nome":"Gilberto",
			"disciplinas": [
				"Orientação a objetos",
				"Programação Web: back-end",
				"Dispositivos Móveis I",
				"Dispositivos Móveis II",
				"Fábrica de Software IV"
			],
			"horarios": {
				"dom":[],
				"seg":[],
				"ter":[1,2,3,4],
				"qua":[1,2,3,4],
				"qui":[1,2,3,4],
				"sex":[1,2,3,4],
				"sab":[]
			}
		},
		{
			"nome":"Marco",
			"disciplinas": [
                "Programação Web: front-end",
				"Metodologia Científica para Computação",
				"Fábrica de Software IV"
			],
			"horarios": {
				"dom":[],
				"seg":[1,2,3,4],
				"ter":[1,2,3,4],
				"qua":[1,2,3,4],
				"qui":[],
				"sex":[],
				"sab":[]
			}
		},
		{
			"nome":"Rosa",
			"disciplinas": [
				"Relações Étnico-Raciais",
				"Metodologia Científica para Computação"
			],
			"horarios": {
				"dom":[],
				"seg":[1,2,3,4],
				"ter":[1,2,3,4],
				"qua":[],
				"qui":[1,2,3,4],
				"sex":[],
				"sab":[]
			}
		},
		{ 
			"nome":"Pinho",
			"disciplinas": [
				"Matemática Computacional"
			],
			"horarios": {
				"dom":[],
				"seg":[],
				"ter":[],
				"qua":[1,2,3,4],
				"qui":[],
				"sex":[],
				"sab":[]
			}
		},
		{ 
			"nome":"Valéria",
			"disciplinas": [
				"Fundamentos em Negócios"
			],
			"horarios": {
				"dom":[],
				"seg":[1,2],
				"ter":[1,2],
				"qua":[],
				"qui":[],
				"sex":[],
				"sab":[]
			}
		}
	]
};