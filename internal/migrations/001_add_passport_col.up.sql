-- Для таблицы train_load_unload_store.loads
ALTER TABLE train_load_unload_store.loads 
ADD COLUMN train_passport integer;

-- Для таблицы train_load_unload_store.unloads
ALTER TABLE train_load_unload_store.unloads 
ADD COLUMN train_passport integer;