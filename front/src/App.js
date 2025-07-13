// src/App.js

import React, { useState, useEffect, useMemo } from 'react';
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
  IconButton,
  createFilterOptions
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import { useForm, SubmitHandler, FormProvider, useFormContext, useFieldArray, Controller } from "react-hook-form"
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

// --- Componente Principal da Aplicação ---
function App() {
//const dadosIniciais = localStorage.getItem('dadosGeracaoHorario') ? JSON.parse(localStorage.getItem('dadosGeracaoHorario')) : dadosPadrao;
    // Use memo para evitar re-renderizações desnecessárias
    const dadosIniciais = React.useMemo(() => {
        const storedData = localStorage.getItem('dadosGeracaoHorario');
        if (storedData) return JSON.parse(storedData);
        else return dadosPadrao; // Retorna os dados padrão se não houver no localStorage
    }, []);

    // 1. Inicializa o react-hook-form com todos os dados iniciais
  // Todos os useState para dados do formulário foram removidos!
  const methods = useForm({
    defaultValues: dadosIniciais
  });
  const { control, register, watch } = methods;

  // --- Estados para a API e exibição do resultado ---
  const [jobId, setJobId] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState(null);
  const [horarioGerado, setHorarioGerado] = useState(null);

  const [professoresCores, setProfessoresCores] = useState({});
  const [professoresDisciplinas, setProfessoresDisciplinas] = useState({});

  const { fields: turmas, append: addTurmas, remove: removeTurmas } = useFieldArray({
    control,
    name: "turmas" // O "caminho" para o array nos dados do formulário
  });

  const { fields: disciplinas, append: addDisciplinas, remove: removeDisciplinas } = useFieldArray({
    control,
    name: "disciplinas" // O "caminho" para o array nos dados do formulário
  });

 const { fields: disciplinas_unidas, append: addDisciplinasUnidas, remove: removeDisciplinasUnidas } = useFieldArray({
    control,
    name: "disciplinas_unidas" // O "caminho" para o array nos dados do formulário
    });

  const { fields: professores, append: addProfessores, remove: removeProfessores } = useFieldArray({
    control,
    name: "professores" // O "caminho" para o array nos dados do formulário
  });

  // --- Função para iniciar a geração do horário ---
  const handleGerarHorario = async (formData) => {
    setIsLoading(true);
    setError(null);
    setHorarioGerado(null);
    setJobId(null);

    
    const professoresCores = {};
    const professoresDisciplinas = {};
    formData.professores.forEach((prof, index) => {
      professoresCores[prof.nome] = PALETA_CORES[index % PALETA_CORES.length];
      prof.disciplinas.forEach(disc => {
        if (!professoresDisciplinas[disc]) { // Evita sobreescrever se outra prof. lecionar a mesma disciplina
          professoresDisciplinas[disc] = prof.nome;
        }
      });
    });
    setProfessoresCores(professoresCores);
    setProfessoresDisciplinas(professoresDisciplinas);  

    try {
      console.log('Dados enviados para a API:', formData);
      // Salvar no local storage:
      localStorage.setItem('dadosGeracaoHorario', JSON.stringify(formData));


      const response = await axios.post(`${API_BASE_URL}/horario/jobs`, formData);
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
  //const turmasOptions = turmas.map(t => t.nome).filter(Boolean);
  //const disciplinasOptions = disciplinas.map(d => d.nome).filter(Boolean);

  // --- NOVO: Adicione estes 'useMemo' para criar os mapas de cores e professores ---
    
  return (
    // 3. Envolvemos toda a aplicação com o FormProvider
    <FormProvider {...methods}>
      {/* 4. Usamos a tag <form> e a função handleSubmit do hook */}
      <form onSubmit={methods.handleSubmit(handleGerarHorario)}>
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
                      variant="outlined"
                      {...register(`turmas.${index}.nome`)} // Caminho para o campo
                    />
                    <IconButton onClick={() => removeTurmas(index)} color="error">
                      <DeleteIcon />
                    </IconButton>
                  </Stack>
                  <Controller
                    control={control}
                    name={`turmas.${index}.horarios`}
                    render={({ field }) => (
                    // O `field` do Controller nos dá { onChange, onBlur, value, ref }
                    <TabelaHorariosCheckbox
                        horarios={field.value}
                        onHorarioChange={field.onChange} // Conectamos o onChange do hook ao nosso componente
                    />
                    )}
                />
                </Box>
              ))}
            </Stack>
            <Button sx={{mt: 2}} onClick={() => addTurmas({ nome: '', horarios: {} })}>
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
                    <TextField label="Nome da Disciplina" fullWidth
                        {...register(`disciplinas.${index}.nome`)}
                    />
                    <IconButton onClick={() => removeDisciplinas(index)} color="error"><DeleteIcon /></IconButton>
                  </Stack>
                    <Controller
                    name={`disciplinas.${index}.turma`}
                    control={control} // Obtido via useFormContext()
                    defaultValue={""} 
                    render={({ field, fieldState: { error } }) => (
                        <Autocomplete
                        options={[]}
                        filterOptions={(options, params) => {
                            return watch("turmas").map(d => d.nome).filter(Boolean);
                        }}
                        value={field.value || []} // Garante que o valor nunca seja undefined
                        onChange={(event, newValue) => {
                            field.onChange(newValue);
                        }}
                        onBlur={field.onBlur} // Informa ao RHF quando o campo perde o foco
                        renderInput={(params) => (
                            <TextField
                            {...params}
                            label="Turma"
                            error={!!error} // Exibe o estado de erro se houver
                            helperText={error?.message} // Exibe a mensagem de erro da validação
                            />
                        )}
                        />
                    )}
                    />
                  <Stack direction="row" spacing={2} sx={{mt: 2}}>
                    <TextField label="Aulas Semanais" type="number" fullWidth
                        {...register(`disciplinas.${index}.aulas`, {
                            valueAsNumber: true
                        })}
                    />
                    <TextField label="Agrupar Aulas" type="number" fullWidth
                        {...register(`disciplinas.${index}.agrupar`, {
                            valueAsNumber: true
                        })}
                    />
                   </Stack>
                </Box>
              ))}
            </Stack>
            <Button sx={{mt: 2}} onClick={() => addDisciplinas({ nome: '', turma: null, aulas: 2, agrupar: 2 })}>
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
                                <TextField label="Nome do Grupo" fullWidth
                                {...register(`disciplinas_unidas.${index}.grupo`)}
                                />
                                <IconButton onClick={() => removeDisciplinasUnidas(index)} color="error"><DeleteIcon /></IconButton>
                            </Stack>
                            <Controller
                            name={`disciplinas_unidas.${index}.disciplinas`}
                            control={control} // Obtido via useFormContext()
                            defaultValue={[]} 
                            render={({ field, fieldState: { error } }) => (
                                <Autocomplete
                                multiple
                                options={[]}
                                filterOptions={(options, params) => {
                                    return watch("disciplinas").map(d => d.nome).filter(Boolean);
                                }}
                                value={field.value || []} // Garante que o valor nunca seja undefined
                                onChange={(event, newValue) => {
                                    field.onChange(newValue);
                                }}
                                onBlur={field.onBlur} // Informa ao RHF quando o campo perde o foco
                                renderInput={(params) => (
                                    <TextField
                                    {...params}
                                    label="Disciplinas para Unir"
                                    error={!!error} // Exibe o estado de erro se houver
                                    helperText={error?.message} // Exibe a mensagem de erro da validação
                                    />
                                )}
                                />
                            )}
                            />
                        </Box>
                    ))}
                </Stack>
                <Button sx={{mt: 2}} onClick={() => addDisciplinasUnidas({ grupo: '', disciplinas: [] })}>
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
                                <TextField label="Nome do Professor" fullWidth 
                                    {...register(`professores.${index}.nome`)}
                                />
                                <IconButton onClick={() => removeProfessores(index)} color="error"><DeleteIcon /></IconButton>
                            </Stack>
                             
                            <Controller
                            name={`professores.${index}.disciplinas`}
                            control={control} // Obtido via useFormContext()
                            defaultValue={[]} 
                            render={({ field, fieldState: { error } }) => (
                                <Autocomplete
                                multiple
                                options={[]}
                                filterOptions={(options, params) => {
                                    return watch("disciplinas").map(d => d.nome).filter(Boolean);
                                }}
                                value={field.value || []} // Garante que o valor nunca seja undefined
                                onChange={(event, newValue) => {
                                    field.onChange(newValue);
                                }}
                                onBlur={field.onBlur} // Informa ao RHF quando o campo perde o foco
                                renderInput={(params) => (
                                    <TextField
                                    {...params}
                                    label="Disciplinas que o professor ministra"
                                    error={!!error} // Exibe o estado de erro se houver
                                    helperText={error?.message} // Exibe a mensagem de erro da validação
                                    />
                                )}
                                />
                            )}
                            />

                             <Typography variant="subtitle2" sx={{mt: 2, mb: 1}}>Preferências de horários do professor:</Typography>
                            <Controller
                                control={control}
                                name={`professores.${index}.horarios`}
                                render={({ field }) => (
                                // O `field` do Controller nos dá { onChange, onBlur, value, ref }
                                <TabelaHorariosCheckbox
                                    horarios={field.value}
                                    onHorarioChange={field.onChange} // Conectamos o onChange do hook ao nosso componente
                                />
                                )}
                            />
                        </Box>
                    ))}
                </Stack>
                <Button sx={{mt: 2}} onClick={() => addProfessores({ nome: '', disciplinas: [], horarios: {} })}>
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
              <Button type="submit" variant="contained" color="primary" size="large" disabled={isLoading}>
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
                                const professor = disciplina ? professoresDisciplinas[disciplina] : null;
                                const cor = professor ? professoresCores[professor] : 'transparent';
                                
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
    </form>
    </FormProvider>
  );
}

export default App;