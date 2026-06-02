'use client';

import { useState, useRef, useCallback, useEffect } from 'react';
import { motion } from 'framer-motion';
import { Search, Link as LinkIcon, X, Loader2, User, AtSign } from 'lucide-react';
import { cn, isValidInstagramUrl, isValidInstagramUsername, extractInstagramShortcode } from '@/lib/utils';

interface SearchBarProps {
  mode: 'link' | 'user';
  onSearch: (query: string) => void;
  isLoading?: boolean;
  autoFocus?: boolean;
}

export default function SearchBar({ mode, onSearch, isLoading, autoFocus }: SearchBarProps) {
  const [query, setQuery] = useState('');
  const [isValid, setIsValid] = useState<boolean | null>(null);
  const [isFocused, setIsFocused] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (autoFocus && inputRef.current) {
      inputRef.current.focus();
    }
  }, [autoFocus]);

  const validateQuery = useCallback((value: string) => {
    if (!value) {
      setIsValid(null);
      return;
    }
    if (mode === 'link') {
      setIsValid(isValidInstagramUrl(value));
    } else {
      setIsValid(isValidInstagramUsername(value));
    }
  }, [mode]);

  // Re-validate when mode changes
  useEffect(() => {
    if (query) validateQuery(query);
  }, [mode, query, validateQuery]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setQuery(value);
    validateQuery(value);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = query.trim();
    if (!trimmed) return;

    if (mode === 'link' && isValidInstagramUrl(trimmed)) {
      onSearch(trimmed);
    } else if (mode === 'user' && isValidInstagramUsername(trimmed)) {
      onSearch(trimmed);
    }
  };

  const handlePaste = async () => {
    try {
      const text = await navigator.clipboard.readText();
      const trimmed = text.trim();
      setQuery(trimmed);
      validateQuery(trimmed);
      if (mode === 'link' && isValidInstagramUrl(trimmed)) {
        onSearch(trimmed);
      } else if (mode === 'user' && isValidInstagramUsername(trimmed)) {
        onSearch(trimmed);
      }
    } catch {
      // Clipboard API not available
    }
  };

  const handleClear = () => {
    setQuery('');
    setIsValid(null);
    inputRef.current?.focus();
  };

  const canSubmit = query.trim() && (
    (mode === 'link' && isValidInstagramUrl(query.trim())) ||
    (mode === 'user' && isValidInstagramUsername(query.trim()))
  );

  const placeholder = mode === 'link'
    ? 'Paste Instagram link here...'
    : 'Enter Instagram username...';

  const hintText = mode === 'link'
    ? 'Please enter a valid Instagram URL (instagram.com/p/..., instagram.com/reel/...)'
    : 'Enter a valid Instagram username (letters, numbers, dots, underscores)';

  const InputIcon = mode === 'link' ? LinkIcon : AtSign;

  return (
    <div className="w-full max-w-2xl mx-auto">
      <form onSubmit={handleSubmit} className="relative">
        <div
          className={cn(
            'relative flex items-center gap-2 sm:gap-3 w-full px-3 sm:px-5 py-3 sm:py-4 rounded-2xl transition-all duration-300',
            'glass',
            isFocused && 'ring-2 ring-brand-500/50 shadow-lg shadow-brand-500/10',
            isValid === false && 'ring-2 ring-red-500/50',
            isLoading && 'opacity-75 pointer-events-none'
          )}
        >
          {/* Icon */}
          <div className="flex-shrink-0">
            {isLoading ? (
              <Loader2 className="h-5 w-5 text-brand-400 animate-spin" />
            ) : (
              <InputIcon className="h-5 w-5 text-surface-400" />
            )}
          </div>

          {/* Input */}
          <input
            ref={inputRef}
            type="text"
            value={query}
            onChange={handleChange}
            onFocus={() => setIsFocused(true)}
            onBlur={() => setIsFocused(false)}
            placeholder={placeholder}
            className="flex-1 bg-transparent text-base text-surface-100 placeholder-surface-500 outline-none border-none min-w-0"
            disabled={isLoading}
            autoComplete="off"
            spellCheck={false}
          />

          {/* Action buttons */}
          <div className="flex items-center gap-2">
            {query && !isLoading && (
              <button
                type="button"
                onClick={handleClear}
                className="flex items-center justify-center p-1.5 rounded-lg text-surface-400 hover:text-surface-100 hover:bg-white/5 transition-colors"
                aria-label="Clear input"
              >
                <X className="h-4 w-4" />
              </button>
            )}

            {!query && !isLoading && (
              <button
                type="button"
                onClick={handlePaste}
                className="hidden sm:flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium text-brand-300 hover:text-surface-100 hover:bg-brand-500/10 transition-colors"
              >
                <LinkIcon className="h-3.5 w-3.5" />
                Paste
              </button>
            )}

            <motion.button
              type="submit"
              disabled={!canSubmit || isLoading}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              className={cn(
                'flex items-center gap-1.5 sm:gap-2 px-3 sm:px-5 py-2 rounded-xl text-xs sm:text-sm font-semibold transition-all duration-200',
                canSubmit
                  ? 'bg-gradient-to-r from-brand-600 to-accent-600 text-white shadow-lg shadow-brand-500/25 hover:shadow-brand-500/40'
                  : 'bg-surface-800 text-surface-400 cursor-not-allowed'
              )}
            >
              {isLoading ? (
                <>
                  <Loader2 className="h-4 w-4 animate-spin" />
                  <span className="hidden sm:inline">Fetching...</span>
                </>
              ) : (
                <>
                  <Search className="h-4 w-4" />
                  <span className="hidden sm:inline">Search</span>
                </>
              )}
            </motion.button>
          </div>
        </div>

        {/* Validation hint */}
        <motion.div
          initial={{ opacity: 0, y: -5 }}
          animate={{
            opacity: isValid === false ? 1 : 0,
            y: isValid === false ? 0 : -5,
          }}
          className="absolute -bottom-6 left-5"
        >
          <span className="text-xs text-red-400">{hintText}</span>
        </motion.div>
      </form>
    </div>
  );
}
