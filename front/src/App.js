// src/App.js

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import {
  Container,
  Typography,
  Box,
  Button,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Checkbox,
  Autocomplete,
  CircularProgress,
  Alert,
  Tabs,
  Tab,
  Card,
  CardContent,
  Stack,
  Divider,
  TextField,
  IconButton
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import { aulas as dadosPadrao } from './aulas';

// --- Configurações da Aplicação ---
const QUANTIDADE_TEMPOS = 4;
const API_BASE_URL = 'http://localhost:8080';
const DIAS_SEMANA_MAP = {
  seg: 'Segunda', ter: 'Terça', qua: 'Quarta', qui: 'Quinta', sex: 'Sexta', sab: 'Sábado', dom: 'Domingo'
};
const DIAS_SEMANA = ['dom', 'seg', 'ter', 'qua', 'qui', 'sex', 'sab'];

const PALETA_CORES = [
  '#649cb1ff', '#f18841ff', '#84f8c8ff', '#9c9b53ff', '#6c76adff', '#ce74f2ff',
  '#179fe9ff', '#ec2d2dff', '#2b8a27ff', '#f79c15ff', '#3339f3ff', '#f17cc0ff',
  '#a75656ff', '#4e958a', '#a01e6aff', '#bec02bff', '#589140ff', '#1f9e78'
];

// Função para obter uma cor de texto contrastante (preto ou branco)
const getContrastColor = (hexColor) => {
  if (!hexColor) return '#000';
  const r = parseInt(hexColor.substr(1, 2), 16);
  const g = parseInt(hexColor.substr(3, 2), 16);
  const b = parseInt(hexColor.substr(5, 2), 16);
  const yiq = ((r * 299) + (g * 587) + (b * 114)) / 1000;
  return (yiq >= 128) ? '#000' : '#fff';
};

// --- Componente Reutilizável para a Tabela de Horários com Checkboxes ---
const TabelaHorariosCheckbox = ({ horarios, onHorarioChange }) => {
  const handleCheckboxChange = (dia, tempo) => {
    const novosHorarios = { ...horarios };
    const temposDoDia = novosHorarios[dia] || [];
    if (temposDoDia.includes(tempo)) {
      novosHorarios[dia] = temposDoDia.filter(t => t !== tempo);
    } else {
      novosHorarios[dia] = [...temposDoDia, tempo].sort((a, b) => a - b);
    }
    onHorarioChange(novosHorarios);
  };

  return (
    <TableContainer component={Paper} variant="outlined" sx={{ mt: 2 }}>
      <Table size="small">
        <TableHead>
          <TableRow>
            <TableCell>Dia</TableCell>
            {Array.from({ length: QUANTIDADE_TEMPOS }, (_, i) => i + 1).map(tempo => (
              <TableCell key={tempo} align="center">Tempo {tempo}</TableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>
          {DIAS_SEMANA.map(dia => (
            <TableRow key={dia}>
              <TableCell>{DIAS_SEMANA_MAP[dia]}</TableCell>
              {Array.from({ length: QUANTIDADE_TEMPOS }, (_, i) => i + 1).map(tempo => (
                <TableCell key={tempo} align="center">
                  <Checkbox
                    checked={horarios[dia]?.includes(tempo) || false}
                    onChange={() => handleCheckboxChange(dia, tempo)}
                  />
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};


// --- Funções Genéricas de Manipulação dos Arrays de Estado ---
const handleStateChange = (setter, index, field, value) => {
setter(prevState => {
    const newState = [...prevState];
    newState[index] = { ...newState[index], [field]: value };
    return newState;
});
};

const handleAddItem = (setter, newItem) => {
setter(prevState => [...prevState, newItem]);
};

const handleRemoveItem = (setter, index) => {
setter(prevState => prevState.filter((_, i) => i !== index));
};

// --- Componente Principal da Aplicação ---
function App() {
//const dadosIniciais = localStorage.getItem('dadosGeracaoHorario') ? JSON.parse(localStorage.getItem('dadosGeracaoHorario')) : dadosPadrao;
    // Use memo para evitar re-renderizações desnecessárias
    const dadosIniciais = React.useMemo(() => {
        const storedData = localStorage.getItem('dadosGeracaoHorario');
        if (storedData) return JSON.parse(storedData);
        else return dadosPadrao; // Retorna os dados padrão se não houver no localStorage
    }, []);

  // --- Estados para os dados da aplicação ---
  const [turmas, setTurmas] = useState(dadosIniciais.turmas);
  const [disciplinas, setDisciplinas] = useState(dadosIniciais.disciplinas);
  const [professores, setProfessores] = useState(dadosIniciais.professores);
  const [disciplinas_unidas, setDisciplinasUnidas] = useState(dadosIniciais.disciplinas_unidas);

  // --- Estados para a API e exibição do resultado ---
  const [jobId, setJobId] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [horarioGerado, setHorarioGerado] = useState(null);

  // --- Função para iniciar a geração do horário ---
  const handleGerarHorario = async () => {
    setIsLoading(true);
    setError(null);
    setHorarioGerado(null);
    setJobId(null);

    const dadosParaApi = {
      turmas,
      disciplinas,
      professores,
      disciplinas_unidas
    };

    try {
      console.log('Dados enviados para a API:', dadosParaApi);
      // Salvar no local storage:
      localStorage.setItem('dadosGeracaoHorario', JSON.stringify(dadosParaApi));


      const response = await axios.post(`${API_BASE_URL}/horario/jobs`, dadosParaApi);
      const { id } = response.data;
      setJobId(id);
    } catch (err) {
      setError('Falha ao iniciar a geração do horário. Verifique se a API está rodando e o CORS está configurado.');
      setIsLoading(false);
    }
  };

  // --- Efeito para verificar o status do Job (Polling) ---
  useEffect(() => {
    if (!jobId) return;

    const intervalId = setInterval(async () => {
      try {
        const response = await axios.get(`${API_BASE_URL}/horario/jobs/${jobId}`);
        const { status, result } = response.data;

        if (status === 'completed') {
          setHorarioGerado(result);
          setIsLoading(false);
          clearInterval(intervalId);
        } else if (status === 'failed') {
          setError('Ocorreu um erro na API ao gerar o horário.');
          setIsLoading(false);
          clearInterval(intervalId);
        }
      } catch (err) {
        setError('Falha ao consultar o status do job.');
        setIsLoading(false);
        clearInterval(intervalId);
      }
    }, 2000); 

    return () => clearInterval(intervalId);
  }, [jobId]);

  // --- Opções para os Autocompletes ---
  const turmasOptions = turmas.map(t => t.nome).filter(Boolean);
  const disciplinasOptions = disciplinas.map(d => d.nome).filter(Boolean);

  // --- NOVO: Adicione estes 'useMemo' para criar os mapas de cores e professores ---
  
  // Cria um mapa para associar cada professor a uma cor
  const professorColorMap = React.useMemo(() => {
    const map = {};
    professores.forEach((prof, index) => {
      map[prof.nome] = PALETA_CORES[index % PALETA_CORES.length];
    });
    return map;
  }, [professores]);

  // Cria um mapa para associar cada disciplina ao seu professor (pegando o primeiro que a leciona)
  const disciplinaParaProfessorMap = React.useMemo(() => {
    const map = {};
    professores.forEach(prof => {
      prof.disciplinas.forEach(disc => {
        if (!map[disc]) { // Evita sobreescrever se outra prof. lecionar a mesma disciplina
          map[disc] = prof.nome;
        }
      });
    });
    return map;
  }, [professores]);
  
  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Typography variant="h3" component="h1" gutterBottom align="center">
        Gerador de Horários ADS
      </Typography>

      <Stack spacing={4}>

        {/* --- Painel de Turmas --- */}
        <Card variant="outlined">
          <CardContent>
            <Typography variant="h5" component="h2" gutterBottom>
              1. Turmas
            </Typography>
            <Stack spacing={3} divider={<Divider />}>
              {turmas.map((turma, index) => (
                <Box key={index}>
                  <Stack direction="row" spacing={2} alignItems="center">
                    <TextField
                      label="Nome da Turma"
                      fullWidth
                      value={turma.nome}
                      onChange={(e) => handleStateChange(setTurmas, index, 'nome', e.target.value)}
                      variant="outlined"
                    />
                    <IconButton onClick={() => handleRemoveItem(setTurmas, index)} color="error">
                      <DeleteIcon />
                    </IconButton>
                  </Stack>
                  <TabelaHorariosCheckbox 
                      horarios={turma.horarios} 
                      onHorarioChange={(h) => handleStateChange(setTurmas, index, 'horarios', h)} 
                  />
                </Box>
              ))}
            </Stack>
            <Button sx={{mt: 2}} onClick={() => handleAddItem(setTurmas, { nome: '', horarios: {} })}>
              + Adicionar Nova Turma
            </Button>
          </CardContent>
        </Card>

        {/* --- Painel de Disciplinas --- */}
        <Card variant="outlined">
          <CardContent>
            <Typography variant="h5" component="h2" gutterBottom>2. Disciplinas</Typography>
            <Stack spacing={3} divider={<Divider />}>
              {disciplinas.map((disciplina, index) => (
                <Box key={index}>
                  <Stack direction="row" spacing={2} alignItems="center" sx={{mb: 2}}>
                    <TextField label="Nome da Disciplina" fullWidth value={disciplina.nome} onChange={(e) => handleStateChange(setDisciplinas, index, 'nome', e.target.value)} />
                    <IconButton onClick={() => handleRemoveItem(setDisciplinas, index)} color="error"><DeleteIcon /></IconButton>
                  </Stack>
                  <Autocomplete
                      options={turmasOptions}
                      value={disciplina.turma || null}
                      onChange={(e, newValue) => handleStateChange(setDisciplinas, index, 'turma', newValue)}
                      renderInput={(params) => <TextField {...params} label="Turma" />}
                  />
                  <Stack direction="row" spacing={2} sx={{mt: 2}}>
                      <TextField label="Aulas Semanais" type="number" fullWidth value={disciplina.aulas} onChange={(e) => handleStateChange(setDisciplinas, index, 'aulas', parseInt(e.target.value, 10) || 0)} />
                      <TextField label="Agrupar Aulas" type="number" fullWidth value={disciplina.agrupar} onChange={(e) => handleStateChange(setDisciplinas, index, 'agrupar', parseInt(e.target.value, 10) || 0)} />
                  </Stack>
                </Box>
              ))}
            </Stack>
            <Button sx={{mt: 2}} onClick={() => handleAddItem(setDisciplinas, { nome: '', turma: null, aulas: 2, agrupar: 2 })}>
              + Adicionar Nova Disciplina
            </Button>
          </CardContent>
        </Card>

        {/* --- Painel de Unir Disciplinas --- */}
        <Card variant="outlined">
            <CardContent>
                <Typography variant="h5" component="h2" gutterBottom>3. Unir Disciplinas</Typography>
                <Stack spacing={3} divider={<Divider />}>
                    {disciplinas_unidas.map((grupo, index) => (
                        <Box key={index}>
                            <Stack direction="row" spacing={2} alignItems="center" sx={{mb: 2}}>
                                <TextField label="Nome do Grupo" fullWidth value={grupo.grupo} onChange={(e) => handleStateChange(setDisciplinasUnidas, index, 'grupo', e.target.value)} />
                                <IconButton onClick={() => handleRemoveItem(setDisciplinasUnidas, index)} color="error"><DeleteIcon /></IconButton>
                            </Stack>
                            <Autocomplete
                                multiple
                                options={disciplinasOptions}
                                value={grupo.disciplinas}
                                onChange={(e, newValue) => handleStateChange(setDisciplinasUnidas, index, 'disciplinas', newValue)}
                                renderInput={(params) => <TextField {...params} label="Disciplinas para Unir" />}
                            />
                        </Box>
                    ))}
                </Stack>
                <Button sx={{mt: 2}} onClick={() => handleAddItem(setDisciplinasUnidas, { grupo: '', disciplinas: [] })}>
                    + Adicionar Novo Grupo
                </Button>
            </CardContent>
        </Card>

        {/* --- Painel de Professores --- */}
        <Card variant="outlined">
            <CardContent>
                <Typography variant="h5" component="h2" gutterBottom>4. Professores</Typography>
                <Stack spacing={3} divider={<Divider />}>
                    {professores.map((professor, index) => (
                        <Box key={index}>
                             <Stack direction="row" spacing={2} alignItems="center" sx={{mb: 2}}>
                                <TextField label="Nome do Professor" fullWidth value={professor.nome} onChange={(e) => handleStateChange(setProfessores, index, 'nome', e.target.value)} />
                                <IconButton onClick={() => handleRemoveItem(setProfessores, index)} color="error"><DeleteIcon /></IconButton>
                            </Stack>
                             <Autocomplete
                                multiple
                                options={disciplinasOptions}
                                value={professor.disciplinas}
                                onChange={(e, newValue) => handleStateChange(setProfessores, index, 'disciplinas', newValue)}
                                renderInput={(params) => <TextField {...params} label="Disciplinas que o professor ministra" />}
                            />
                             <Typography variant="subtitle2" sx={{mt: 2, mb: 1}}>Preferências de horários do professor:</Typography>
                            <TabelaHorariosCheckbox 
                                horarios={professor.horarios}
                                onHorarioChange={(h) => handleStateChange(setProfessores, index, 'horarios', h)}
                            />
                        </Box>
                    ))}
                </Stack>
                <Button sx={{mt: 2}} onClick={() => handleAddItem(setProfessores, { nome: '', disciplinas: [], horarios: {} })}>
                    + Adicionar Novo Professor
                </Button>
            </CardContent>
        </Card>

        {/* --- Painel de Geração e Resultado --- */}
        <Card variant="outlined">
            <CardContent>
            <Typography variant="h5" component="h2" gutterBottom>
              5. Gerar e Visualizar Horário
            </Typography>
            <Divider sx={{ my: 2 }}/>
            <Box sx={{ display: 'flex', justifyContent: 'center', my: 3 }}>
              <Button variant="contained" color="primary" size="large" onClick={handleGerarHorario} disabled={isLoading}>
                {isLoading ? <CircularProgress size={24} color="inherit" /> : 'Gerar Horário'}
              </Button>
            </Box>

            {error && <Alert severity="error" sx={{ mt: 2 }}>{error}</Alert>}
            
            {isLoading && !horarioGerado && (
                <Box sx={{textAlign: 'center', mt: 2}}>
                    <Typography>Processando... O horário está sendo gerado.</Typography>
                    <Typography variant="caption">Job ID: {jobId}</Typography>
                </Box>
            )}

{horarioGerado && (
          <Box sx={{ mt: 4 }}>
            <Typography variant="h6" gutterBottom>Horário Gerado com Sucesso!</Typography>
            
            <Stack spacing={4} sx={{mt: 3}}>
              {horarioGerado.map((resultado, index) => {
                return (
                  <Box key={index}>
                    <Typography variant="h6" component="h3" sx={{mb: 1}}>
                       {resultado.turma}
                    </Typography>
                    <TableContainer component={Paper} variant="outlined">
                      <Table stickyHeader size="small">
                        <TableHead>
                          <TableRow>
                            {DIAS_SEMANA.map(dia => (
                              <TableCell key={dia} align="center" sx={{fontWeight: 'bold', backgroundColor: '#e0e0e0'}}>
                                {DIAS_SEMANA_MAP[dia]}
                              </TableCell>
                            ))}
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {Array.from({ length: QUANTIDADE_TEMPOS }, (_, i) => i).map(indexTempo => (
                            <TableRow key={indexTempo}>
                              {DIAS_SEMANA.map(dia => {
                                const horarioDaTurma = resultado.horario;
                                const diaInfo = horarioDaTurma?.find(h => h.dia === dia);
                                const disciplina = diaInfo?.tempos[indexTempo];
                                const professor = disciplina ? disciplinaParaProfessorMap[disciplina] : null;
                                const cor = professor ? professorColorMap[professor] : 'transparent';
                                
                                return (
                                  <TableCell
                                    key={dia}
                                    align="center"
                                    sx={{
                                      backgroundColor: cor,
                                      border: '1px solid #ccc',
                                      color: getContrastColor(cor),
                                      fontWeight: '500',
                                      minWidth: '120px'
                                    }}
                                  >
                                    {disciplina || '---'}
                                  </TableCell>
                                );
                              })}
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </Box>
                )
              })}
            </Stack>
          </Box>
        )}
         </CardContent>
        </Card>
      </Stack>
    </Container>
  );
}

export default App;