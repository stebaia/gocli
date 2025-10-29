package generator

import (
	"fmt"
	"fline-cli/internal/ui"
	"fline-cli/internal/utils"
	"strings"
)

// ModelGenerator generates models, services, repositories and BLoCs from JSON
type ModelGenerator struct {
	modelName   string
	jsonData    map[string]interface{}
	endpoint    string
	packageName string
	writer      *utils.FileWriter
	naming      *utils.NamingHelper
	logger      *ui.Logger
}

// NewModelGenerator creates a new model generator
func NewModelGenerator(
	modelName string,
	jsonData map[string]interface{},
	endpoint string,
	packageName string,
	writer *utils.FileWriter,
) *ModelGenerator {
	return &ModelGenerator{
		modelName:   modelName,
		jsonData:    jsonData,
		endpoint:    endpoint,
		packageName: packageName,
		writer:      writer,
		naming:      utils.NewNamingHelper(modelName),
		logger:      ui.NewLogger("model"),
	}
}

// Generate generates all components
func (g *ModelGenerator) Generate() error {
	if err := g.generateModel(); err != nil {
		return fmt.Errorf("failed to generate model: %w", err)
	}

	if err := g.generateService(); err != nil {
		return fmt.Errorf("failed to generate service: %w", err)
	}

	if err := g.generateRepository(); err != nil {
		return fmt.Errorf("failed to generate repository: %w", err)
	}

	if err := g.generateBloc(); err != nil {
		return fmt.Errorf("failed to generate bloc: %w", err)
	}

	return nil
}

func (g *ModelGenerator) generateModel() error {
	fields := g.generateFields()
	constructorParams := g.generateConstructorParams()
	props := g.generateProps()

	content := fmt.Sprintf(`import 'package:json_annotation/json_annotation.dart';
import 'package:equatable/equatable.dart';

part '%s.g.dart';

@JsonSerializable()
class %s extends Equatable {
  %s

  const %s({
    %s
  });

  factory %s.fromJson(Map<String, dynamic> json) =>
      _$%sFromJson(json);

  Map<String, dynamic> toJson() => _$%sToJson(this);

  @override
  List<Object?> get props => [
    %s
  ];
}
`,
		g.naming.SnakeCase(),
		g.naming.PascalCase(),
		fields,
		g.naming.PascalCase(),
		constructorParams,
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		props,
	)

	return g.writer.WriteFile(
		fmt.Sprintf("lib/model/%s.dart", g.naming.SnakeCase()),
		content,
	)
}

func (g *ModelGenerator) generateService() error {
	content := fmt.Sprintf(`import 'package:dio/dio.dart';
import 'package:retrofit/retrofit.dart';
import 'package:%s/model/%s.dart';

part '%s_service.g.dart';

@RestApi()
abstract class %sService {
  factory %sService(Dio dio) = _%sService;

  @GET('%s')
  Future<List<%s>> getAll();

  @GET('%s/{id}')
  Future<%s> getById(@Path('id') String id);

  @POST('%s')
  Future<%s> create(@Body() %s %s);

  @PUT('%s/{id}')
  Future<%s> update(
    @Path('id') String id,
    @Body() %s %s,
  );

  @DELETE('%s/{id}')
  Future<void> delete(@Path('id') String id);
}
`,
		g.packageName,
		g.naming.SnakeCase(),
		g.naming.SnakeCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.endpoint,
		g.naming.PascalCase(),
		g.endpoint,
		g.naming.PascalCase(),
		g.endpoint,
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.CamelCase(),
		g.endpoint,
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.CamelCase(),
		g.endpoint,
	)

	return g.writer.WriteFile(
		fmt.Sprintf("lib/network/service/%s_service.dart", g.naming.SnakeCase()),
		content,
	)
}

func (g *ModelGenerator) generateRepository() error {
	content := fmt.Sprintf(`import 'package:logger/logger.dart';
import 'package:%s/model/%s.dart';
import 'package:%s/network/service/%s_service.dart';

class %sRepository {
  final %sService _service;
  final Logger _logger;

  %sRepository({
    required %sService service,
    required Logger logger,
  })  : _service = service,
        _logger = logger;

  Future<List<%s>> getAll() async {
    try {
      return await _service.getAll();
    } catch (e) {
      _logger.e('Error fetching %ss', error: e);
      rethrow;
    }
  }

  Future<%s> getById(String id) async {
    try {
      return await _service.getById(id);
    } catch (e) {
      _logger.e('Error fetching %s', error: e);
      rethrow;
    }
  }

  Future<%s> create(%s %s) async {
    try {
      return await _service.create(%s);
    } catch (e) {
      _logger.e('Error creating %s', error: e);
      rethrow;
    }
  }

  Future<%s> update(String id, %s %s) async {
    try {
      return await _service.update(id, %s);
    } catch (e) {
      _logger.e('Error updating %s', error: e);
      rethrow;
    }
  }

  Future<void> delete(String id) async {
    try {
      await _service.delete(id);
    } catch (e) {
      _logger.e('Error deleting %s', error: e);
      rethrow;
    }
  }
}
`,
		g.packageName,
		g.naming.SnakeCase(),
		g.packageName,
		g.naming.SnakeCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.CamelCase(),
		g.naming.CamelCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.CamelCase(),
		g.naming.CamelCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
	)

	return g.writer.WriteFile(
		fmt.Sprintf("lib/repositories/%s_repository.dart", g.naming.SnakeCase()),
		content,
	)
}

func (g *ModelGenerator) generateBloc() error {
	// Generate bloc
	blocContent := g.generateBlocFile()
	if err := g.writer.WriteFile(
		fmt.Sprintf("lib/state_management/bloc/%s/%s_bloc.dart",
			g.naming.SnakeCase(), g.naming.SnakeCase()),
		blocContent,
	); err != nil {
		return err
	}

	// Generate event
	eventContent := g.generateEventFile()
	if err := g.writer.WriteFile(
		fmt.Sprintf("lib/state_management/bloc/%s/%s_event.dart",
			g.naming.SnakeCase(), g.naming.SnakeCase()),
		eventContent,
	); err != nil {
		return err
	}

	// Generate state
	stateContent := g.generateStateFile()
	return g.writer.WriteFile(
		fmt.Sprintf("lib/state_management/bloc/%s/%s_state.dart",
			g.naming.SnakeCase(), g.naming.SnakeCase()),
		stateContent,
	)
}

func (g *ModelGenerator) generateBlocFile() string {
	return fmt.Sprintf(`import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:equatable/equatable.dart';
import 'package:%s/repositories/%s_repository.dart';
import 'package:%s/model/%s.dart';

part '%s_event.dart';
part '%s_state.dart';

class %sBloc extends Bloc<%sEvent, %sState> {
  final %sRepository _repository;

  %sBloc({required %sRepository repository})
      : _repository = repository,
        super(%sInitial()) {
    on<Fetch%ss>(_onFetch%ss);
    on<Fetch%s>(_onFetch%s);
    on<Create%s>(_onCreate%s);
    on<Update%s>(_onUpdate%s);
    on<Delete%s>(_onDelete%s);
  }

  Future<void> _onFetch%ss(
    Fetch%ss event,
    Emitter<%sState> emit,
  ) async {
    emit(%sLoading());
    try {
      final items = await _repository.getAll();
      emit(%ssLoaded(items));
    } catch (e) {
      emit(%sError(e.toString()));
    }
  }

  Future<void> _onFetch%s(
    Fetch%s event,
    Emitter<%sState> emit,
  ) async {
    emit(%sLoading());
    try {
      final item = await _repository.getById(event.id);
      emit(%sLoaded(item));
    } catch (e) {
      emit(%sError(e.toString()));
    }
  }

  Future<void> _onCreate%s(
    Create%s event,
    Emitter<%sState> emit,
  ) async {
    emit(%sLoading());
    try {
      await _repository.create(event.%s);
      emit(%sCreated());
    } catch (e) {
      emit(%sError(e.toString()));
    }
  }

  Future<void> _onUpdate%s(
    Update%s event,
    Emitter<%sState> emit,
  ) async {
    emit(%sLoading());
    try {
      await _repository.update(event.id, event.%s);
      emit(%sUpdated());
    } catch (e) {
      emit(%sError(e.toString()));
    }
  }

  Future<void> _onDelete%s(
    Delete%s event,
    Emitter<%sState> emit,
  ) async {
    emit(%sLoading());
    try {
      await _repository.delete(event.id);
      emit(%sDeleted());
    } catch (e) {
      emit(%sError(e.toString()));
    }
  }
}
`,
		g.packageName, g.naming.SnakeCase(),
		g.packageName, g.naming.SnakeCase(),
		g.naming.SnakeCase(), g.naming.SnakeCase(),
		g.naming.PascalCase(), g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.CamelCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.CamelCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
	)
}

func (g *ModelGenerator) generateEventFile() string {
	return fmt.Sprintf(`part of '%s_bloc.dart';

abstract class %sEvent extends Equatable {
  const %sEvent();

  @override
  List<Object> get props => [];
}

class Fetch%ss extends %sEvent {}

class Fetch%s extends %sEvent {
  final String id;

  const Fetch%s(this.id);

  @override
  List<Object> get props => [id];
}

class Create%s extends %sEvent {
  final %s %s;

  const Create%s(this.%s);

  @override
  List<Object> get props => [%s];
}

class Update%s extends %sEvent {
  final String id;
  final %s %s;

  const Update%s(this.id, this.%s);

  @override
  List<Object> get props => [id, %s];
}

class Delete%s extends %sEvent {
  final String id;

  const Delete%s(this.id);

  @override
  List<Object> get props => [id];
}
`,
		g.naming.SnakeCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.CamelCase(),
		g.naming.PascalCase(), g.naming.CamelCase(),
		g.naming.CamelCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.CamelCase(),
		g.naming.PascalCase(), g.naming.CamelCase(),
		g.naming.CamelCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
	)
}

func (g *ModelGenerator) generateStateFile() string {
	return fmt.Sprintf(`part of '%s_bloc.dart';

abstract class %sState extends Equatable {
  const %sState();

  @override
  List<Object> get props => [];
}

class %sInitial extends %sState {}

class %sLoading extends %sState {}

class %ssLoaded extends %sState {
  final List<%s> items;

  const %ssLoaded(this.items);

  @override
  List<Object> get props => [items];
}

class %sLoaded extends %sState {
  final %s item;

  const %sLoaded(this.item);

  @override
  List<Object> get props => [item];
}

class %sCreated extends %sState {}

class %sUpdated extends %sState {}

class %sDeleted extends %sState {}

class %sError extends %sState {
  final String message;

  const %sError(this.message);

  @override
  List<Object> get props => [message];
}
`,
		g.naming.SnakeCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(), g.naming.PascalCase(),
		g.naming.PascalCase(),
	)
}

func (g *ModelGenerator) generateFields() string {
	var fields []string
	for key, value := range g.jsonData {
		fieldType := g.getFieldType(value)
		fields = append(fields, fmt.Sprintf("final %s %s;", fieldType, key))
	}
	return strings.Join(fields, "\n  ")
}

func (g *ModelGenerator) generateConstructorParams() string {
	var params []string
	for key := range g.jsonData {
		params = append(params, fmt.Sprintf("required this.%s", key))
	}
	return strings.Join(params, ",\n    ")
}

func (g *ModelGenerator) generateProps() string {
	var props []string
	for key := range g.jsonData {
		props = append(props, key)
	}
	return strings.Join(props, ",\n    ")
}

func (g *ModelGenerator) getFieldType(value interface{}) string {
	if value == nil {
		return "dynamic"
	}

	switch v := value.(type) {
	case bool:
		return "bool"
	case float64:
		// Check if it's actually an int
		if v == float64(int64(v)) {
			return "int"
		}
		return "double"
	case string:
		return "String"
	case []interface{}:
		if len(v) == 0 {
			return "List<dynamic>"
		}
		elemType := g.getFieldType(v[0])
		return fmt.Sprintf("List<%s>", elemType)
	case map[string]interface{}:
		return "Map<String, dynamic>"
	default:
		return "dynamic"
	}
}
